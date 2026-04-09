package handler

import (
	"encoding/json"
	"net/http"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Signup godoc
// @Summary Sign up a new user
// @Description Register a user by email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authRequest true "Signup payload"
// @Success 201 {object} data.User
// @Failure 400 {string} string "invalid json"
// @Failure 400 {string} string "email and password required"
// @Failure 500 {string} string "failed to signup"
// @Router /signup [post]
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	var user data.User
	user, err := h.authService.Signup(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "failed to signup", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate user credentials and return access/refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authRequest true "Login payload"
// @Success 200 {object} service.AuthTokens
// @Failure 400 {string} string "invalid json"
// @Failure 401 {string} string "invalid credentials"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Exchange a valid refresh token for a new token pair
// @Tags auth
// @Accept json
// @Produce json
// @Param request body refreshRequest true "Refresh token payload"
// @Success 200 {object} service.AuthTokens
// @Failure 400 {string} string "invalid json"
// @Failure 401 {string} string "invalid refresh token"
// @Router /refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

// Logout godoc
// @Summary Log out a user
// @Description Revoke the provided refresh token
// @Tags auth
// @Accept json
// @Param request body refreshRequest true "Refresh token payload"
// @Success 204 "No Content"
// @Failure 400 {string} string "invalid json"
// @Failure 500 {string} string "failed to logout"
// @Router /logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
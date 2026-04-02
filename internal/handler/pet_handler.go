package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/errs"
	"github.com/andreishemetov/pawpal/internal/middleware"
	"github.com/andreishemetov/pawpal/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type PetHandler struct {
	service service.PetStore
}

func NewPetHandler(service service.PetStore) *PetHandler {
	return &PetHandler{service: service}
}

// GetPets godoc
// @Summary List pets
// @Description Get pets with pagination and filtering
// @Tags pets
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param type query string false "Pet type"
// @Param q query string false "Search by name"
// @Param sort query string false "Sort field"
// @Success 200 {object} data.PetsListResponse
// @Failure 500 {string} string "failed to load pets"
// @Router /pets [get]
func (h *PetHandler) GetPets(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserId(w, r)
	if !ok {
		return
	}

	userRole, ok := getUserRole(w, r)
	if !ok {
		return
	}

	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)

	page = clamp(page, 1, 1_000_000)
	limit = clamp(limit, 1, 100)
	petType := r.URL.Query().Get("type") // optional
	q := r.URL.Query().Get("q")          // optional search
	sort := r.URL.Query().Get("sort")

	pets, total, err := h.service.GetAll(r.Context(), userID, userRole, data.PetQuery{
		Page:  page,
		Limit: limit,
		Type:  petType,
		Q:     q,
		Sort:  sort,
	}) // optional

	if err != nil {
		http.Error(w, "failed to load pets", http.StatusInternalServerError)
		return
	}

	resp := data.PetsListResponse{
		Items: pets,
	}
	resp.Meta.Page = page
	resp.Meta.Limit = limit
	resp.Meta.Total = total

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreatePet godoc
// @Summary Create a pet
// @Description Create a new pet
// @Tags pets
// @Accept json
// @Produce json
// @Param pet body data.Pet true "Pet payload"
// @Success 201 {object} data.Pet
// @Failure 405 {string} string "Method not allowed"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 400 {string} string "name is required"
// @Failure 500 {string} string "failed to create pet"
// @Router /pets [post]
func (h *PetHandler) PostPet(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserId(w, r)
	if !ok {
		return
	}

	reqID := chiMiddleware.GetReqID(r.Context()) // string
	fmt.Printf("Request ID: %s", reqID)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var pet data.Pet

	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validation
	if pet.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	pet.UserID = userID

	created, err := h.service.Add(r.Context(), pet)
	if err != nil {
		http.Error(w, "failed to create pet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetCountPets godoc
// @Summary Get count of pets
// @Tags pets
// @Success 200 {object} data.CountResponse
// @Failure 500 {string} string "failed to load pets"
// @Router /pets/count [get]
func (h *PetHandler) GetCountPets(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserId(w, r)
	if !ok {
		return
	}

	userRole, ok := getUserRole(w, r)
	if !ok {
		return
	}

	_, total, err := h.service.GetAll(r.Context(), userID, userRole, data.PetQuery{})
	if err != nil {
		http.Error(w, "failed to load pets", http.StatusInternalServerError)
		return
	}
	count := total
	w.Header().Set("Content-Type", "application/json")
	response := data.CountResponse{
		Count: count,
	}

	json.NewEncoder(w).Encode(response)
}

// GetPetByID godoc
// @Summary Get pet by ID
// @Tags pets
// @Produce json
// @Param id path int true "Pet ID"
// @Success 200 {object} data.Pet
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "pet not found"
// @Failure 500 {string} string "internal error"
// @Router /pets/{id} [get]
func (h *PetHandler) GetPetByID(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserId(w, r)
	if !ok {
		return
	}

	userRole, ok := getUserRole(w, r)
	if !ok {
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	pet, err := h.service.GetByID(r.Context(), userID, userRole, id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, "pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pet)
}

// DeletePetByID godoc
// @Summary Delete pet by ID
// @Tags pets
// @Param id path int true "Pet ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "pet not found"
// @Failure 500 {string} string "internal error"
// @Router /pets/{id} [delete]
func (h *PetHandler) DeletePetByID(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserId(w, r)
	if !ok {
		return
	}

	userRole, ok := getUserRole(w, r)
	if !ok {
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	isDeleted, err := h.service.DeleteByID(r.Context(), userID, userRole, id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !isDeleted {
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePet godoc
// @Summary Update pet by ID
// @Tags pets
// @Param id path int true "Pet ID"
// @Success 200 {object} data.Pet
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Router /pets/{id} [put]
func (h *PetHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserId(w, r)
	if !ok {
		return
	}

	userRole, ok := getUserRole(w, r)
	if !ok {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.ErrorResponse{
			Error: "invalid id",
		})
		return
	}

	var pet data.Pet
	err = json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.ErrorResponse{
			Error: "invalid json",
		})
		return
	}

	pet, err = h.service.Update(r.Context(), userID, userRole, id, pet)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errs.ErrorResponse{
				Error: "pet not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errs.ErrorResponse{
			Error: "internal error",
		})
		return
	}

	json.NewEncoder(w).Encode(pet)
}

func getUserId(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return 0, false
	}
	fmt.Println("current user id:", id)
	return id, true
}

func getUserRole(w http.ResponseWriter, r *http.Request) (string, bool) {
	role, ok := middleware.GetUserRole(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return "", false
	}
	fmt.Println("current user role:", role)
	return role, true
}

func parseIntQuery(r *http.Request, key string, def int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return def
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return n
}

func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

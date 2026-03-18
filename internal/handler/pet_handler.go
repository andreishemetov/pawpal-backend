package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/errs"
	"github.com/andreishemetov/pawpal/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CountResponse struct {
	Count int `json:"count"`
}

type PetHandler struct {
	service service.PetStore
}

func NewPetHandler(service service.PetStore) *PetHandler {
	return &PetHandler{service: service}
}

func (h *PetHandler) GetPets(w http.ResponseWriter, r *http.Request) {
	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)

	page = clamp(page, 1, 1_000_000)
	limit = clamp(limit, 1, 100)
	petType := r.URL.Query().Get("type") // optional
	q := r.URL.Query().Get("q")          // optional search
	sort := r.URL.Query().Get("sort")

	pets, total, err := h.service.GetAll(r.Context(), data.PetQuery{
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

	resp := map[string]any{
		"items": pets,
		"meta": map[string]any{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *PetHandler) PostPet(w http.ResponseWriter, r *http.Request) {
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

	h.service.Add(r.Context(), pet)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pet)
}

func (h *PetHandler) GetCountPets(w http.ResponseWriter, r *http.Request) {
	_, total, err := h.service.GetAll(r.Context(), data.PetQuery{})
	if err != nil {
		http.Error(w, "failed to load pets", http.StatusInternalServerError)
		return
	}
	count := total
	w.Header().Set("Content-Type", "application/json")
	response := CountResponse{
		Count: count,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *PetHandler) GetPetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	pet, err := h.service.GetByID(r.Context(), id)
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

func (h *PetHandler) DeletePetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	isDeleted, err := h.service.DeleteByID(r.Context(), id)
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

func (h *PetHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "invalid id",
		})
		return
	}

	var pet data.Pet
	err = json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "invalid json",
		})
		return
	}

	pet, err = h.service.Update(r.Context(), id, pet)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: "pet not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "internal error",
		})
		return
	}

	json.NewEncoder(w).Encode(pet)
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

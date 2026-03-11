package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/repo"
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

	w.Header().Set("Content-Type", "application/json")

	pets, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to load pets", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pets)
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
	pets, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to load pets", http.StatusInternalServerError)
		return
	}
	count := len(pets)
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
		if errors.Is(err, repo.ErrNotFound) {
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

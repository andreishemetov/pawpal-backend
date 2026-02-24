package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/service"
	"github.com/go-chi/chi/v5"
)

var pets = []data.Pet{}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CountResponse struct {
	Count int `json:"count"`
}

type PetHandler struct {
	service *service.PetService
}

func NewPetHandler(service *service.PetService) *PetHandler {
	return &PetHandler{service: service}
}

func (h *PetHandler) GetPets(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(h.service.GetAll())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func (h *PetHandler) PostPet(w http.ResponseWriter, r *http.Request) {
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

	h.service.Add(pet)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pet)
}

func (h *PetHandler) GetCountPets(w http.ResponseWriter, r *http.Request) {
	count := len(h.service.GetAll())
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

	pet, found := h.service.GetByID(id)
	if found {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pet)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: "pet not found",
	})
}

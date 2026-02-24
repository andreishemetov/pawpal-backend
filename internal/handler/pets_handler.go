package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/go-chi/chi/v5"
)

var pets = []data.Pet{}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CountResponse struct {
	Count int `json:"count"`
}

func GetPets(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(pets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func PostPet(w http.ResponseWriter, r *http.Request) {
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

	pets = append(pets, pet)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pet)
}

func GetCountPets(w http.ResponseWriter, r *http.Request) {
	count := len(pets)
	w.Header().Set("Content-Type", "application/json")
	response := CountResponse{
		Count: count,
	}

	json.NewEncoder(w).Encode(response)
}

func GetPetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	for _, pet := range pets {
		if pet.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(pet)
			return
		}
	}

	// If not found
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: "pet not found",
	})
}

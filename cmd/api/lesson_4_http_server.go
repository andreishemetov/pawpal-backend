package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/andreishemetov/pawpal/internal/data"
)

func lesson4() {

	fmt.Println("Lesson 4 starting...")

	fmt.Println("Server running on :8080")

	var pets = []data.Pet{}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/pets/count", func(w http.ResponseWriter, r *http.Request) {
		countPetsHandler(w, r, &pets)
	})
	http.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getPets(w, r, &pets)
			return
		}
		if r.Method == http.MethodPost {
			createPet(w, r, &pets)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getPets(w http.ResponseWriter, r *http.Request, pets *[]data.Pet) {

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(pets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func createPet(w http.ResponseWriter, r *http.Request, pets *[]data.Pet) {
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

	*pets = append(*pets, pet)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pet)
}

func countPetsHandler(w http.ResponseWriter, r *http.Request, pets *[]data.Pet) {
	count := len(*pets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]int{
		"count": count,
	}
	json.NewEncoder(w).Encode(response)
}

/*
curl -X POST http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"Charlie","age":3,"visits":0}'


curl http://localhost:8080/pets

*/

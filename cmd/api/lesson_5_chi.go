package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/andreishemetov/pawpal/internal/handler"
)

func lesson5() {

	fmt.Println("Lesson 5 starting...")

	fmt.Println("Server running on :8080")

	router := chi.NewRouter()

	router.Get("/health", getHealth)
	router.Get("/pets", handler.GetPets)
	router.Post("/pets", handler.PostPet)
	router.Get("/pets/{id}", handler.GetPetByID)
	router.Get("/pets/count", handler.GetCountPets)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}


/*
curl -X POST http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"Charlie","age":3,"visits":0}'


curl http://localhost:8080/pets

*/
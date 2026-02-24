package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andreishemetov/pawpal/internal/handler"
	"github.com/andreishemetov/pawpal/internal/service"
	"github.com/go-chi/chi/v5"
)

/*
HTTP layer  →  Handler
Handler     →  Service
Service     →  Data model
*/

func lesson5() {

	fmt.Println("Lesson 5 starting...")

	fmt.Println("Server running on :8080")

	router := chi.NewRouter()
	petService := service.NewPetService()
	petHandler := handler.NewPetHandler(petService)

	router.Get("/health", getHealth)
	router.Get("/pets", petHandler.GetPets)
	router.Post("/pets", petHandler.PostPet)
	router.Get("/pets/{id}", petHandler.GetPetByID)
	router.Get("/pets/count", petHandler.GetCountPets)

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

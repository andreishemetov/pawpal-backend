package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andreishemetov/pawpal/internal/handler"
	"github.com/andreishemetov/pawpal/internal/middleware"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/andreishemetov/pawpal/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Standard useful middlewares
	router.Use(chiMiddleware.RequestID)                // generates request IDs
	router.Use(chiMiddleware.RealIP)                   // uses X-Forwarded-For, etc.
	router.Use(chiMiddleware.Timeout(5 * time.Second)) // request timeout
	router.Use(middleware.Logging)                     // our custom logger
	router.Use(chiMiddleware.Recoverer)                // recover panics

	dsn := "postgres://pawpal:pawpal_pass@localhost:5432/pawpal_dev?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repo.NewUserPostgresRepo(db)
	authService := service.NewAuthService(userRepo, "very-secret-key")
	authHandler := handler.NewAuthHandler(authService)

	router.Post("/signup", authHandler.Signup)
	router.Post("/login", authHandler.Login)

	authMiddleware := middleware.AuthMiddleware("very-secret-key")

	petRepo := repo.NewPetPostgresRepo(db)
	petHandler := handler.NewPetHandler(petRepo)

	router.Get("/health", getHealth)

	router.Group(
		func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/pets", petHandler.GetPets)
			r.Post("/pets", petHandler.PostPet)
			r.Get("/pets/{id}", petHandler.GetPetByID)
			r.Get("/pets/count", petHandler.GetCountPets)
			r.Delete("/pets/{id}", petHandler.DeletePetByID)
			r.Put("/pets/{id}", petHandler.UpdatePet)
		},
	)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

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

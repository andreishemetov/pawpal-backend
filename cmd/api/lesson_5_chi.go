package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andreishemetov/pawpal/internal/handler"
	"github.com/andreishemetov/pawpal/internal/middleware"
	"github.com/andreishemetov/pawpal/internal/notify"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/andreishemetov/pawpal/internal/service"
	"github.com/andreishemetov/pawpal/internal/worker"
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
	router.Use(chiMiddleware.RequestID)                 // generates request IDs
	router.Use(chiMiddleware.RealIP)                    // uses X-Forwarded-For, etc.
	router.Use(chiMiddleware.Timeout(15 * time.Second)) // request timeout
	router.Use(middleware.Logging)                      // our custom logger
	router.Use(chiMiddleware.Recoverer)                 // recover panics

	dsn := "postgres://pawpal:pawpal_pass@localhost:5432/pawpal_dev?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	reminderRepo := repo.NewReminderPostgresRepo(db)
	userRepo := repo.NewUserPostgresRepo(db)
	refreshTokenRepo := repo.NewRefreshTokenPostgresRepo(db)

	reminderHandler := handler.NewReminderHandler(reminderRepo)
	logNotifier := notify.NewLogNotifier()

	reminderWorker := worker.NewReminderWorker(reminderRepo, userRepo, logNotifier)
	go reminderWorker.Start(context.Background())

	authService := service.NewAuthService(userRepo, refreshTokenRepo, "very-secret-key")
	authHandler := handler.NewAuthHandler(authService)

	router.Post("/signup", authHandler.Signup)
	router.Post("/login", authHandler.Login)
	router.Post("/refresh", authHandler.Refresh)
	router.Post("/logout", authHandler.Logout)

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
			r.Post("/reminders", reminderHandler.CreateReminder)
		},
	)

	router.Group(
		func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(middleware.RequireRole("admin"))
			r.Get("/admin/ping", getAdminPing)
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

func getAdminPing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("admin ok"))
}

/*
curl -X POST http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"id":1,"name":"Charlie","age":3,"visits":0}'

curl -X POST http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"Charlie","age":3,"visits":0}'


curl http://localhost:8080/pets

curl http://localhost:8080/pets \
  -H "Authorization: Bearer YOUR_TOKEN"

curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

  curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

  docker compose exec db psql -U pawpal -d pawpal_dev -c "
UPDATE users
SET role = 'admin'
WHERE email = 'admin@example.com';


curl -X POST http://localhost:8080/reminders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "pet_id": 1,
    "message": "Give medicine",
    "remind_at": "2026-04-02T12:00:00Z"
  }'

*/

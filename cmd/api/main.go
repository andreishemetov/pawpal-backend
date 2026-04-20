// @title PawPal API
// @version 1.0
// @description API for managing pets in PawPal.
// @termsOfService https://example.com/terms/

// @contact.name API Support
// @contact.url https://example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/andreishemetov/pawpal/docs"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andreishemetov/pawpal/internal/cache"
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

func main() {

	fmt.Println("Lesson 5 starting...")

	fmt.Println("Server running on :8080")

	router := chi.NewRouter()

	// Standard useful middlewares
	router.Use(chiMiddleware.RequestID)                 // generates request IDs
	router.Use(chiMiddleware.RealIP)                    // uses X-Forwarded-For, etc.
	router.Use(chiMiddleware.Timeout(15 * time.Second)) // request timeout
	router.Use(middleware.Logging)                      // our custom logger
	router.Use(chiMiddleware.Recoverer)                 // recover panics

	// dsn := "postgres://pawpal:pawpal_pass@localhost:5432/pawpal_dev?sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is required")
	}
	petCache := cache.NewRedisCache(redisAddr, 30*time.Second)

	reminderRepo := repo.NewReminderPostgresRepo(db)
	userRepo := repo.NewUserPostgresRepo(db)
	refreshTokenRepo := repo.NewRefreshTokenPostgresRepo(db)

	reminderHandler := handler.NewReminderHandler(reminderRepo)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	router.Post("/signup", authHandler.Signup)
	router.Post("/login", authHandler.Login)
	router.Post("/refresh", authHandler.Refresh)
	router.Post("/logout", authHandler.Logout)

	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	petRepo := repo.NewPetPostgresRepo(db)
	petHandler := handler.NewPetHandler(petRepo, petCache)

	rateLimitMiddleware := middleware.RateLimitMiddleware(petCache, 10, 1*time.Minute)

	router.Get("/health", getHealth)

	router.Group(
		func(r chi.Router) {
			r.Use(authMiddleware)
			r.Use(rateLimitMiddleware)
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

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Println("Server running on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
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

swag init -g cmd/api/main.go -o docs

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



docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0001_init.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0002_add_indexes.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0003_users.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0004_add_user_id_to_pets.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0005_refresh_tokens.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0006_add_user_role.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0007_reminders.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0008_reminder_delivery_fields.up.sql

*/

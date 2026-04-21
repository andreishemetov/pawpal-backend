package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/andreishemetov/pawpal/docs"
	"github.com/andreishemetov/pawpal/internal/cache"
	"github.com/andreishemetov/pawpal/internal/config"
	"github.com/andreishemetov/pawpal/internal/handler"
	"github.com/andreishemetov/pawpal/internal/middleware"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/andreishemetov/pawpal/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RunAPI wires dependencies and serves HTTP until the process exits.
func RunAPI(ctx context.Context, cfg *config.Config) error {
	db, err := OpenPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer func() { _ = db.Close() }()

	petCache := cache.NewRedisCache(cfg.RedisAddr, cfg.PetCacheTTL)

	reminderRepo := repo.NewReminderPostgresRepo(db)
	userRepo := repo.NewUserPostgresRepo(db)
	refreshTokenRepo := repo.NewRefreshTokenPostgresRepo(db)

	reminderHandler := handler.NewReminderHandler(reminderRepo)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	router := chi.NewRouter()

	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Timeout(cfg.HTTPTimeout))
	router.Use(middleware.Logging)
	router.Use(chiMiddleware.Recoverer)

	router.Post("/signup", authHandler.Signup)
	router.Post("/login", authHandler.Login)
	router.Post("/refresh", authHandler.Refresh)
	router.Post("/logout", authHandler.Logout)

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	petRepo := repo.NewPetPostgresRepo(db)
	petHandler := handler.NewPetHandler(petRepo, petCache)

	rateLimitMiddleware := middleware.RateLimitMiddleware(petCache, cfg.RateLimitMax, cfg.RateLimitWindow)

	router.Get("/health", healthHandler)

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
			r.Get("/admin/ping", adminPingHandler)
		},
	)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	log.Println("Server listening on", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		return err
	}
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func adminPingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("admin ok"))
}

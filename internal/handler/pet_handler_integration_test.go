package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/andreishemetov/pawpal/internal/cache"
	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/middleware"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestServer(t *testing.T) (*httptest.Server, *sql.Tx, *sql.DB) {
	dsn := "postgres://pawpal:pawpal_pass@localhost:5432/pawpal_dev?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(
		`INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)`,
		101, "handler-integration@example.com", "test-hash", "user",
	); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	repository := repo.NewPetPostgresRepo(tx)
	handler := NewPetHandler(repository, cache.NewRedisCache("localhost:6379", time.Minute))
	router := chi.NewRouter()
	router.Use(middleware.AuthMiddleware("test-secret"))
	router.Post("/pets", handler.PostPet)
	router.Get("/pets/{id}", handler.GetPetByID)

	server := httptest.NewServer(router)

	return server, tx, db
}

func makeTestToken(t *testing.T, userID int, role string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  float64(userID),
		"role": role,
	})

	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}

	return signed
}

func TestCreateAndGetPet(t *testing.T) {
	server, tx, db := setupTestServer(t)
	defer server.Close()
	defer tx.Rollback()
	defer db.Close()

	// Create pet
	body := `{"name":"Charlie","type":"Dog","age":3,"visits":0}`

	req, err := http.NewRequest(http.MethodPost, server.URL+"/pets", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+makeTestToken(t, 101, "user"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var created data.Pet
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created pet: %v", err)
	}

	// Get pet by actual ID from response
	req2, err := http.NewRequest(http.MethodGet, server.URL+"/pets/"+strconv.Itoa(created.ID), nil)
	if err != nil {
		t.Fatal(err)
	}
	req2.Header.Set("Authorization", "Bearer "+makeTestToken(t, 101, "user"))

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}
}

func TestCreatePet_Validation(t *testing.T) {
	server, tx, db := setupTestServer(t)
	defer server.Close()
	defer tx.Rollback()
	defer db.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/pets", strings.NewReader(`{"age":3}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+makeTestToken(t, 101, "user"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

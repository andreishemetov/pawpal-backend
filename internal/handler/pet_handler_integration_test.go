package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/repo"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestServer(t *testing.T) (*httptest.Server, *sql.DB) {
	dsn := "postgres://pawpal:pawpal_pass@localhost:5432/pawpal_dev?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}

	repository := repo.NewPetPostgresRepo(db)
	handler := NewPetHandler(repository)
	router := chi.NewRouter()
	router.Post("/pets", handler.PostPet)
	router.Get("/pets/{id}", handler.GetPetByID)

	server := httptest.NewServer(router)

	return server, db
}

func cleanDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM pets")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateAndGetPet(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()

	cleanDB(t, db)

	// Create pet
	body := `{"name":"Charlie","type":"Dog","age":3,"visits":0}`

	resp, err := http.Post(server.URL+"/pets", "application/json", strings.NewReader(body))
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
	resp2, err := http.Get(server.URL + "/pets/" + strconv.Itoa(created.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}
}

func TestCreatePet_Validation(t *testing.T) {
	server, db := setupTestServer(t)
	defer server.Close()

	cleanDB(t, db)

	resp, err := http.Post(server.URL+"/pets", "application/json", strings.NewReader(`{"age":3}`))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

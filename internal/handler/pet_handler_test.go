package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/handler"
	"github.com/andreishemetov/pawpal/internal/service"
	"github.com/go-chi/chi/v5"
)

func setupRouter() (*chi.Mux, *service.PetService) {
	svc := service.NewPetService()
	h := handler.NewPetHandler(svc)

	r := chi.NewRouter()
	r.Get("/pets", h.GetPets)
	r.Post("/pets", h.PostPet)
	r.Get("/pets/{id}", h.GetPetByID)
	return r, svc
}

func TestCreatePetHandler_Success(t *testing.T) {
	r, _ := setupRouter()

	newPet := data.Pet{ID: 1, Name: "Charlie", Age: 3}
	body, _ := json.Marshal(newPet)

	req := httptest.NewRequest(http.MethodPost, "/pets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp data.Pet
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Name != "Charlie" || resp.ID != 1 {
		t.Fatalf("unexpected response: %#v", resp)
	}
}

func TestCreatePetHandler_ValidateName(t *testing.T) {
	r, _ := setupRouter()

	// missing name
	payload := map[string]interface{}{"id": 2, "age": 2}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/pets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when name missing, got %d, body: %s", rec.Code, rec.Body.String())
	}
}

func TestGetPetByID_NotFoundAndFound(t *testing.T) {
	r, svc := setupRouter()

	// seed one pet
	svc.Add(data.Pet{ID: 10, Name: "Milo", Age: 4})

	// not found
	req1 := httptest.NewRequest(http.MethodGet, "/pets/999", nil)
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d, body: %s", rec1.Code, rec1.Body.String())
	}

	// found
	req2 := httptest.NewRequest(http.MethodGet, "/pets/10", nil)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec2.Code, rec2.Body.String())
	}
	var got data.Pet
	if err := json.NewDecoder(rec2.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.ID != 10 || got.Name != "Milo" {
		t.Fatalf("unexpected pet: %#v", got)
	}
}

/*
# run all tests
go test ./...

# run with verbose output for only service tests
go test ./internal/service -v

# detect race conditions
go test ./... -race

# run a single test:
go test ./internal/handler -run TestCreatePetHandler_Success -v
*/

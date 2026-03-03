package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/handler"
	"github.com/go-chi/chi/v5"
)

// fakeStore implements service.PetStore for tests
type fakeStore struct {
	pets []data.Pet
}

func (f *fakeStore) GetAll() []data.Pet {
	// return copy
	return append([]data.Pet(nil), f.pets...)
}

func (f *fakeStore) Add(p data.Pet) {
	f.pets = append(f.pets, p)
}

func (f *fakeStore) GetByID(id int) (*data.Pet, bool) {
	for i := range f.pets {
		if f.pets[i].ID == id {
			return &f.pets[i], true
		}
	}
	return nil, false
}

func (f *fakeStore) DeleteById(id int) bool {
	for i := range f.pets {
		if f.pets[i].ID == id {
			f.pets = append(f.pets[:i], f.pets[i+1:]...)
			return true
		}
	}
	return false
}

// helper to create router with injected fake
func setupRouterWithFake(fake *fakeStore) *chi.Mux {
	h := handler.NewPetHandler(fake)

	r := chi.NewRouter()
	r.Get("/pets", h.GetPets)
	r.Post("/pets", h.PostPet)
	r.Get("/pets/{id}", h.GetPetByID)
	return r
}

func TestCreatePet_WithFake(t *testing.T) {
	f := &fakeStore{pets: []data.Pet{}}
	r := setupRouterWithFake(f)

	newPet := data.Pet{ID: 1, Name: "Charlie", Age: 3}
	body, _ := json.Marshal(newPet)

	req := httptest.NewRequest(http.MethodPost, "/pets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// Inspect fake's state to ensure Add was called
	if len(f.pets) != 1 {
		t.Fatalf("expected fake to have 1 pet, got %d", len(f.pets))
	}
	if f.pets[0].Name != "Charlie" {
		t.Fatalf("expected name Charlie, got %s", f.pets[0].Name)
	}
}

func TestGetPetByID_WithFake(t *testing.T) {
	f := &fakeStore{pets: []data.Pet{
		{ID: 10, Name: "Milo", Age: 4},
	}}
	r := setupRouterWithFake(f)

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

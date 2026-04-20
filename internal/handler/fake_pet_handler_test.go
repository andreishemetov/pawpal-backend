package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andreishemetov/pawpal/internal/cache"
	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/errs"
	"github.com/andreishemetov/pawpal/internal/handler"
	"github.com/andreishemetov/pawpal/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/go-chi/chi/v5"
)

// fakeStore implements service.PetStore for tests
type fakeStore struct {
	pets []data.Pet
}

func (f *fakeStore) GetAll(ctx context.Context, userID int, userRole string, q data.PetQuery) ([]data.Pet, int, error) {
	pets := append([]data.Pet(nil), f.pets...)
	return pets, len(pets), nil
}

func (f *fakeStore) Add(ctx context.Context, p data.Pet) (data.Pet, error) {
	f.pets = append(f.pets, p)
	return p, nil
}

func (f *fakeStore) GetByID(ctx context.Context, userID int, userRole string, id int) (data.Pet, error) {
	for i := range f.pets {
		if f.pets[i].ID == id {
			return f.pets[i], nil
		}
	}
	return data.Pet{}, errs.ErrNotFound
}

func (f *fakeStore) DeleteByID(ctx context.Context, userID int, userRole string, id int) (bool, error) {
	for i := range f.pets {
		if f.pets[i].ID == id {
			f.pets = append(f.pets[:i], f.pets[i+1:]...)
			return true, nil
		}
	}
	return false, errs.ErrNotFound
}

func (f *fakeStore) Update(ctx context.Context, userID int, userRole string, id int, p data.Pet) (data.Pet, error) {
	for i := range f.pets {
		if f.pets[i].ID == id {
			f.pets[i] = p
			f.pets[i].ID = id
			return f.pets[i], nil
		}
	}
	return data.Pet{}, errs.ErrNotFound
}

// helper to create router with injected fake
func setupRouterWithFake(fake *fakeStore) *chi.Mux {
	h := handler.NewPetHandler(fake, cache.NewRedisCache("localhost:6379", time.Minute))

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware("test-secret"))
	r.Get("/pets", h.GetPets)
	r.Post("/pets", h.PostPet)
	r.Get("/pets/{id}", h.GetPetByID)
	r.Delete("/pets/{id}", h.DeletePetByID)
	return r
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

func TestCreatePet_WithFake(t *testing.T) {
	f := &fakeStore{pets: []data.Pet{}}
	r := setupRouterWithFake(f)

	newPet := data.Pet{ID: 1, Name: "Charlie", Age: 3}
	body, _ := json.Marshal(newPet)

	req := httptest.NewRequest(http.MethodPost, "/pets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+makeTestToken(t, 101, "user"))
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
	req1.Header.Set("Authorization", "Bearer "+makeTestToken(t, 101, "user"))
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d, body: %s", rec1.Code, rec1.Body.String())
	}

	// found
	req2 := httptest.NewRequest(http.MethodGet, "/pets/10", nil)
	req2.Header.Set("Authorization", "Bearer "+makeTestToken(t, 101, "user"))
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

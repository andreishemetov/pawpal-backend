package service_test

import (
	"context"
	"testing"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/service"
)

func TestPetService_AddAndGet(t *testing.T) {
	s := service.NewPetService()

	p1 := data.Pet{ID: 1, Name: "Charlie", Age: 3}
	p2 := data.Pet{ID: 2, Name: "Milo", Age: 2}

	s.Add(context.Background(), p1)
	s.Add(context.Background(), p2)

	all, _ := s.GetAll(context.Background())
	if len(all) != 2 {
		t.Fatalf("expected 2 pets, got %d", len(all))
	}

	got, err := s.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected to find pet ID 1")
	}
	if got.Name != "Charlie" {
		t.Fatalf("expected name Charlie, got %s", got.Name)
	}
}

func TestPetService_GetByID_NotFound(t *testing.T) {
	s := service.NewPetService()

	_, err := s.GetByID(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected not found")
	}
}

func TestPetService_DeleteByID(t *testing.T) {
	// If you implemented DeleteByID(id int) bool
	s := service.NewPetService()
	s.Add(context.Background(), data.Pet{ID: 1, Name: "A"})
	s.Add(context.Background(), data.Pet{ID: 2, Name: "B"})

	deleted, err := s.DeleteByID(context.Background(), 1)
	if !deleted {
		t.Fatalf("expected delete true")
	}

	_, err = s.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected pet 1 to be gone")
	}
}

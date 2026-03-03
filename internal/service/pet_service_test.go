package service_test

import (
	"testing"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/service"
)

func TestPetService_AddAndGet(t *testing.T) {
	s := service.NewPetService()

	p1 := data.Pet{ID: 1, Name: "Charlie", Age: 3}
	p2 := data.Pet{ID: 2, Name: "Milo", Age: 2}

	s.Add(p1)
	s.Add(p2)

	all := s.GetAll()
	if len(all) != 2 {
		t.Fatalf("expected 2 pets, got %d", len(all))
	}

	got, ok := s.GetByID(1)
	if !ok {
		t.Fatalf("expected to find pet ID 1")
	}
	if got.Name != "Charlie" {
		t.Fatalf("expected name Charlie, got %s", got.Name)
	}
}

func TestPetService_GetByID_NotFound(t *testing.T) {
	s := service.NewPetService()

	_, ok := s.GetByID(999)
	if ok {
		t.Fatalf("expected not found")
	}
}

func TestPetService_DeleteByID(t *testing.T) {
	// If you implemented DeleteByID(id int) bool
	s := service.NewPetService()
	s.Add(data.Pet{ID: 1, Name: "A"})
	s.Add(data.Pet{ID: 2, Name: "B"})

	deleted := s.DeleteById(1)
	if !deleted {
		t.Fatalf("expected delete true")
	}

	_, ok := s.GetByID(1)
	if ok {
		t.Fatalf("expected pet 1 to be gone")
	}

	all := s.GetAll()
	if len(all) != 1 {
		t.Fatalf("expected 1 remaining pet, got %d", len(all))
	}
}
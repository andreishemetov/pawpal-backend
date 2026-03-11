package service

import (
	"context"
	"sync"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/repo"
)

type PetService struct {
	mutex sync.RWMutex
	pets  []data.Pet
}

func NewPetService() *PetService {
	return &PetService{
		pets: []data.Pet{},
	}
}

func (s *PetService) GetAll(ctx context.Context) ([]data.Pet, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return append([]data.Pet{}, s.pets...), nil
}

func (s *PetService) Add(ctx context.Context, pet data.Pet) (data.Pet, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.pets = append(s.pets, pet)
	return pet, nil
}

func (s *PetService) GetByID(ctx context.Context, id int) (data.Pet, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for i := range s.pets {
		if s.pets[i].ID == id {
			return s.pets[i], nil
		}
	}
	return data.Pet{}, repo.ErrNotFound
}

func (s *PetService) DeleteByID(ctx context.Context, id int) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := range s.pets {
		if s.pets[i].ID == id {
			s.pets = append(s.pets[:i], s.pets[i+1:]...)
			return true, nil
		}
	}
	return false, repo.ErrNotFound
}

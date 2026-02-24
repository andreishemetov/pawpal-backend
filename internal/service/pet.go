package service

import "github.com/andreishemetov/pawpal/internal/data"

type PetService struct {
	pets []data.Pet
}

func NewPetService() *PetService {
	return &PetService{
		pets: []data.Pet{},
	}
}

func (s *PetService) GetAll() []data.Pet {
	return s.pets
}

func (s *PetService) Add(pet data.Pet) {
	s.pets = append(s.pets, pet)
}

func (s *PetService) GetByID(id int) (*data.Pet, bool) {
	for i := range s.pets {
		if s.pets[i].ID == id {
			return &s.pets[i], true
		}
	}
	return nil, false
}
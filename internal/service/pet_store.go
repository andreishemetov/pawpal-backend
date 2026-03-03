package service

import "github.com/andreishemetov/pawpal/internal/data"

// PetStore describes operations the handler depends on.
// PetService (concrete) will implicitly implement this.
type PetStore interface {
	GetAll() []data.Pet
	Add(p data.Pet)
	GetByID(id int) (*data.Pet, bool)
	DeleteById(id int) bool
}
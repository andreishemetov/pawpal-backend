package service

import (
	"context"

	"github.com/andreishemetov/pawpal/internal/data"
)

// PetStore describes operations the handler depends on.
// PetService (concrete) will implicitly implement this.
type PetStore interface {
	GetAll(ctx context.Context) ([]data.Pet, error)
	Add(ctx context.Context, p data.Pet) (data.Pet, error)
	GetByID(ctx context.Context, id int) (*data.Pet, bool)
	DeleteById(ctx context.Context, id int) bool
}

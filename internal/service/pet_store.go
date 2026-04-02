package service

import (
	"context"

	"github.com/andreishemetov/pawpal/internal/data"
)

// PetStore describes operations the handler depends on.
// PetService (concrete) will implicitly implement this.
type PetStore interface {
	GetAll(ctx context.Context, userID int, userRole string, q data.PetQuery) ([]data.Pet, int, error)
	Add(ctx context.Context, p data.Pet) (data.Pet, error)
	GetByID(ctx context.Context, userID int, userRole string, id int) (data.Pet, error)
	DeleteByID(ctx context.Context, userID int, userRole string, id int) (bool, error)
	Update(ctx context.Context, userID int, userRole string, id int, p data.Pet) (data.Pet, error)
}

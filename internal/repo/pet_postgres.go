package repo

import (
	"context"
	"database/sql"

	"github.com/andreishemetov/pawpal/internal/data"
)

type PetPostgresRepo struct {
	db *sql.DB
}

func NewPetPostgresRepo(db *sql.DB) *PetPostgresRepo {
	return &PetPostgresRepo{db: db}
}

func (r *PetPostgresRepo) GetAll(ctx context.Context) ([]data.Pet, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, type, age, visits FROM pets ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []data.Pet
	for rows.Next() {
		var p data.Pet
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Age, &p.Visits); err != nil {
			return nil, err
		}
		pets = append(pets, p)
	}
	return pets, rows.Err()
}

func (r *PetPostgresRepo) Add(ctx context.Context, p data.Pet) (data.Pet, error) {
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO pets (name, type, age, visits) VALUES ($1,$2,$3,$4) RETURNING id`,
		p.Name, p.Type, p.Age, p.Visits,
	).Scan(&p.ID)
	if err != nil {
		return data.Pet{}, err
	}
	return p, nil
}
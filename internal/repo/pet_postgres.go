package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/andreishemetov/pawpal/internal/errs"
)

type querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type PetPostgresRepo struct {
	q querier
}

func NewPetPostgresRepo(q querier) *PetPostgresRepo {
	return &PetPostgresRepo{q: q}
}

func (r *PetPostgresRepo) GetAll(ctx context.Context, userID int, userRole string, q data.PetQuery) ([]data.Pet, int, error) {
	// sort whitelist (avoid SQL injection)
	orderBy := "id"
	switch q.Sort {
	case "", "id":
		orderBy = "id"
	case "name":
		orderBy = "name"
	case "age":
		orderBy = "age"
	case "visits":
		orderBy = "visits"
	}

	where := []string{"TRUE"}
	args := []any{}
	argN := 1

	if userRole != "admin" {
		where = append(where, fmt.Sprintf("user_id = $%d", argN))
		args = append(args, userID)
		argN++
	}

	if q.Type != "" {
		where = append(where, fmt.Sprintf("type = $%d", argN))
		args = append(args, q.Type)
		argN++
	}

	if q.Q != "" {
		where = append(where, fmt.Sprintf("name ILIKE $%d", argN))
		args = append(args, "%"+q.Q+"%")
		argN++
	}

	whereSQL := strings.Join(where, " AND ")

	// total count
	var total int
	countSQL := "SELECT COUNT(*) FROM pets WHERE " + whereSQL
	if err := r.q.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.Limit

	// data query
	args = append(args, q.Limit, offset)
	limitArg := argN
	offsetArg := argN + 1

	dataSQL := fmt.Sprintf(`
		SELECT id, name, type, age, visits
		FROM pets
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereSQL, orderBy, limitArg, offsetArg)

	rows, err := r.q.QueryContext(ctx, dataSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pets []data.Pet
	for rows.Next() {
		var p data.Pet
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Age, &p.Visits); err != nil {
			return nil, 0, err
		}
		pets = append(pets, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return pets, total, nil
}

func (r *PetPostgresRepo) Add(ctx context.Context, p data.Pet) (data.Pet, error) {
	err := r.q.QueryRowContext(
		ctx,
		`INSERT INTO pets (name, type, age, visits, user_id)
		 VALUES ($1,$2,$3,$4,$5)
		 RETURNING id`,
		p.Name, p.Type, p.Age, p.Visits, p.UserID,
	).Scan(&p.ID)
	if err != nil {
		return data.Pet{}, err
	}
	return p, nil
}

func (r *PetPostgresRepo) GetByID(ctx context.Context, userID int, userRole string, id int) (data.Pet, error) {
	var p data.Pet
	var err error

	if userRole == "admin" {
		err = r.q.QueryRowContext(ctx, `
			SELECT id, name, type, age, visits, user_id
			FROM pets
			WHERE id = $1
		`, id).Scan(&p.ID, &p.Name, &p.Type, &p.Age, &p.Visits, &p.UserID)
	} else {
		err = r.q.QueryRowContext(ctx, `
			SELECT id, name, type, age, visits, user_id
			FROM pets
			WHERE id = $1 AND user_id = $2
		`, id, userID).Scan(&p.ID, &p.Name, &p.Type, &p.Age, &p.Visits, &p.UserID)
	}

	if err == sql.ErrNoRows {
		return data.Pet{}, errs.ErrNotFound
	}

	if err != nil {
		return data.Pet{}, err
	}

	return p, nil
}

func (r *PetPostgresRepo) DeleteByID(ctx context.Context, userID int, userRole string, id int) (bool, error) {
	var res sql.Result
	var err error
	if userRole == "admin" {
		res, err = r.q.ExecContext(ctx,
			`DELETE FROM pets WHERE id = $1`,
			id,
		)
	} else {
		res, err = r.q.ExecContext(ctx,
			`DELETE FROM pets WHERE id = $1 AND user_id = $2`,
			id, userID,
		)
	}
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *PetPostgresRepo) Update(ctx context.Context, userID int, userRole string, id int, p data.Pet) (data.Pet, error) {
	var err error

	if userRole == "admin" {
		err = r.q.QueryRowContext(
			ctx,
			`UPDATE pets
			 SET name = $1, type = $2, age = $3, visits = $4
			 WHERE id = $5
			 RETURNING id, name, type, age, visits, user_id`,
			p.Name, p.Type, p.Age, p.Visits, id,
		).Scan(&p.ID, &p.Name, &p.Type, &p.Age, &p.Visits, &p.UserID)
	} else {
		err = r.q.QueryRowContext(
			ctx,
			`UPDATE pets
			 SET name = $1, type = $2, age = $3, visits = $4
			 WHERE id = $5 AND user_id = $6
			 RETURNING id, name, type, age, visits, user_id`,
			p.Name, p.Type, p.Age, p.Visits, id, userID,
		).Scan(&p.ID, &p.Name, &p.Type, &p.Age, &p.Visits, &p.UserID)
	}

	if err == sql.ErrNoRows {
		return data.Pet{}, errs.ErrNotFound
	}
	if err != nil {
		return data.Pet{}, err
	}
	return p, nil
}

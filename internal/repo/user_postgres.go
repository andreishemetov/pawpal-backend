package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/andreishemetov/pawpal/internal/data"
)

var ErrUserNotFound = errors.New("user not found")

type UserPostgresRepo struct {
	db *sql.DB
}

func NewUserPostgresRepo(db *sql.DB) *UserPostgresRepo {
	return &UserPostgresRepo{db: db}
}

func (r *UserPostgresRepo) Create(ctx context.Context, email, passwordHash string) (data.User, error) {
	var u data.User

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, role
	`, email, passwordHash, "user").
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)

	return u, err
}

func (r *UserPostgresRepo) GetByEmail(ctx context.Context, email string) (data.User, error) {
	var u data.User

	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)

	if err == sql.ErrNoRows {	
		return data.User{}, ErrUserNotFound
	}

	return u, err
}

func (r *UserPostgresRepo) GetByID(ctx context.Context, id int) (data.User, error) {
	var u data.User

	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, role
		FROM users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)

	if err == sql.ErrNoRows {
		return data.User{}, ErrUserNotFound
	}

	return u, err
}

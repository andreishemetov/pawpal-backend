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
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash
	`, email, passwordHash).
		Scan(&u.ID, &u.Email, &u.PasswordHash)

	return u, err
}

func (r *UserPostgresRepo) GetByEmail(ctx context.Context, email string) (data.User, error) {
	var u data.User

	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash)

	if err == sql.ErrNoRows {
		return data.User{}, ErrUserNotFound
	}

	return u, err
}
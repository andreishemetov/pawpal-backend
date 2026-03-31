package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/andreishemetov/pawpal/internal/data"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshTokenPostgresRepo struct {
	db *sql.DB
}

func NewRefreshTokenPostgresRepo(db *sql.DB) *RefreshTokenPostgresRepo {
	return &RefreshTokenPostgresRepo{db: db}
}

func (r *RefreshTokenPostgresRepo) Create(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expiresAt)

	return err
}

func (r *RefreshTokenPostgresRepo) GetValid(ctx context.Context, token string) (data.RefreshToken, error) {
	var rt data.RefreshToken

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE token = $1
		  AND revoked_at IS NULL
		  AND expires_at > now()
	`, token).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
		&rt.RevokedAt,
	)

	if err == sql.ErrNoRows {
		return data.RefreshToken{}, ErrRefreshTokenNotFound
	}

	return rt, err
}

func (r *RefreshTokenPostgresRepo) Revoke(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE token = $1
		  AND revoked_at IS NULL
	`, token)

	return err
}
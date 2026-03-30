package service

import (
	"context"
	"errors"
	"time"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserStore interface {
	Create(ctx context.Context, email, passwordHash string) (data.User, error)
	GetByEmail(ctx context.Context, email string) (data.User, error)
}

type AuthService struct {
	users     UserStore
	jwtSecret []byte
}

func NewAuthService(users UserStore, secret string) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (data.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return data.User{}, err
	}

	return s.users.Create(ctx, email, string(hash))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSecret)
}
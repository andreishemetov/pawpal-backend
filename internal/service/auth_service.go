package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/andreishemetov/pawpal/internal/data"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserStore interface {
	Create(ctx context.Context, email, passwordHash string) (data.User, error)
	GetByEmail(ctx context.Context, email string) (data.User, error)
	GetByID(ctx context.Context, id int) (data.User, error)
}

type RefreshTokenStore interface {
	Create(ctx context.Context, userID int, token string, expiresAt time.Time) error
	GetValid(ctx context.Context, token string) (data.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
}

type AuthService struct {
	users         UserStore
	refreshTokens RefreshTokenStore
	jwtSecret     []byte
}

func NewAuthService(users UserStore, refreshTokens RefreshTokenStore, secret string) *AuthService {
	return &AuthService{
		users:         users,
		refreshTokens: refreshTokens,
		jwtSecret:     []byte(secret),
	}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (data.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return data.User{}, err
	}

	return s.users.Create(ctx, email, string(hash))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (AuthTokens, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return AuthTokens{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return AuthTokens{}, ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return AuthTokens{}, err
	}

	refreshToken, err := generateRandomToken()
	if err != nil {
		return AuthTokens{}, err
	}

	err = s.refreshTokens.Create(ctx, user.ID, refreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.refreshTokens.Revoke(ctx, refreshToken)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (AuthTokens, error) {
	rt, err := s.refreshTokens.GetValid(ctx, refreshToken)
	if err != nil {
		return AuthTokens{}, ErrInvalidCredentials
	}

	user, err := s.users.GetByID(ctx, rt.UserID)
	if err != nil {
		return AuthTokens{}, ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return AuthTokens{}, err
	}

	newRefreshToken, err := generateRandomToken()
	if err != nil {
		return AuthTokens{}, err
	}

	err = s.refreshTokens.Revoke(ctx, refreshToken)
	if err != nil {
		return AuthTokens{}, err
	}

	err = s.refreshTokens.Create(ctx, user.ID, newRefreshToken, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) generateAccessToken(user data.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

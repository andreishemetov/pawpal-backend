package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"
const userRoleKey contextKey = "userRole"

func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(userRoleKey).(string)
	return role, ok
}

// RequireRole allows the request only if the authenticated user has the given role.
// Use after AuthMiddleware so GetUserRole sees JWT claims in context.
func RequireRole(required string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := GetUserRole(r.Context())
			if !ok || role != required {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"]
			if !ok {
				http.Error(w, "missing subject", http.StatusUnauthorized)
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "invalid role", http.StatusUnauthorized)
				return
			}

			var userID int
			switch v := sub.(type) {
			case float64:
				userID = int(v)
			default:
				http.Error(w, "invalid subject", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, userRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/andreishemetov/pawpal/internal/cache"
)

func RateLimitMiddleware(cache *cache.RedisCache, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID, ok := GetUserID(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			key := "rate:user:" + strconv.Itoa(userID)

			count, err := cache.Increment(r.Context(), key, window)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if count > int64(limit) {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
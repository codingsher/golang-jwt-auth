// middleware/jwt.go

package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	jwtLib "github.com/golang-jwt/jwt/v5"
	utilsJwt "github.com/codingsher/user-jwt-auth/internal/utils/jwt" // alias to avoid name collision
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type contextKey string

const UserContextKey contextKey = "userEmail"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Use utilsJwt.CustomClaims for parsing
		token, err := jwtLib.ParseWithClaims(tokenStr, &utilsJwt.CustomClaims{}, func(token *jwtLib.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*utilsJwt.CustomClaims)
		if !ok || claims.Type != "access" {
			http.Error(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


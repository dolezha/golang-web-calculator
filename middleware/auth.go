package middleware

import (
	"calculator/services"
	"calculator/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.RespondWithJSON(w, map[string]string{"error": "Missing authorization token"}, http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				utils.RespondWithJSON(w, map[string]string{"error": "Invalid token format"}, http.StatusUnauthorized)
				return
			}

			claims, err := authService.ValidateToken(tokenParts[1])
			if err != nil {
				utils.RespondWithJSON(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(r *http.Request) (*services.Claims, bool) {
	claims, ok := r.Context().Value(UserContextKey).(*services.Claims)
	return claims, ok
}

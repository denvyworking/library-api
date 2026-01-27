package api

import (
	"context"
	"net/http"
	"strings"
)

func (api *api) RightAuth(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing Authorization header", http.StatusUnauthorized)
		return false
	}

	// Ожидаем формат: "Bearer <token>"
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
		return false
	}

	tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)
	if tokenStr == "" {
		http.Error(w, "empty token", http.StatusUnauthorized)
		return false
	}

	claims, err := api.jwtService.ParseToken(tokenStr)
	if err != nil {
		api.logger.Error("Invalid JWT token", "error", err, "token", tokenStr)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return false
	}

	ctx := context.WithValue(r.Context(), "user_claims", claims)
	*r = *r.WithContext(ctx)

	return true
}

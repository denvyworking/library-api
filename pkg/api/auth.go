package api

import (
	"context"
	"net/http"
	"strings"
)

func (api *api) RightAuth(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	api.logger.Info("Auth header", "header", authHeader)

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
		return false
	}

	tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)
	api.logger.Info("Token to parse", "token", tokenStr)

	claims, err := api.jwtService.ParseToken(tokenStr)
	if err != nil {
		api.logger.Error("JWT parse error", "error", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return false
	}

	ctx := context.WithValue(r.Context(), "user_claims", claims)
	*r = *r.WithContext(ctx)
	return true
}

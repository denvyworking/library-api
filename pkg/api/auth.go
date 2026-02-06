package api

import (
	"context"
	"leti/pkg/auth"
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

// Типизированный ключ для контекста (избегаем коллизий со строковыми ключами)
type contextKey string

const userClaimsKey contextKey = "user_claims"

// GetClaimsFromContext извлекает claims из контекста
func GetClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*auth.Claims)
	return claims, ok
}

func (api *api) RightAuth(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	api.logger.Info("Auth header", "header", authHeader)

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

	ctx := context.WithValue(r.Context(), userClaimsKey, claims)
	*r = *r.WithContext(ctx)
	return true
}

// RequireRole проверяет, что пользователь имеет требуемую роль
func (api *api) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.Role != role {
				api.logger.Warn("Access denied", "required_role", role, "user_role", claims.Role, "user_id", claims.UserID)
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

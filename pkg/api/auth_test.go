package api

import (
	"context"
	"leti/pkg/auth"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetClaimsFromContext(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		claims := &auth.Claims{
			UserID: 123,
			Role:   "admin",
		}

		ctx := context.WithValue(context.Background(), userClaimsKey, claims)
		retrievedClaims, ok := GetClaimsFromContext(ctx)

		require.True(t, ok)
		require.Equal(t, 123, retrievedClaims.UserID)
		require.Equal(t, "admin", retrievedClaims.Role)
	})

	t.Run("not found", func(t *testing.T) {
		ctx := context.Background()
		_, ok := GetClaimsFromContext(ctx)

		require.False(t, ok)
	})

	t.Run("wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userClaimsKey, "not-claims")
		_, ok := GetClaimsFromContext(ctx)

		require.False(t, ok)
	})
}

func TestRequireRole(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	jwtService := auth.NewJWTService("test-secret")

	// Создаем тестовый апи для middleware
	apiInst := &api{
		logger:     logger,
		jwtService: jwtService,
	}

	t.Run("success with correct role", func(t *testing.T) {
		claims := &auth.Claims{
			UserID: 1,
			Role:   "admin",
		}

		ctx := context.WithValue(context.Background(), userClaimsKey, claims)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := apiInst.RequireRole("admin")
		middleware(handler).ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("forbidden with wrong role", func(t *testing.T) {
		claims := &auth.Claims{
			UserID: 1,
			Role:   "user",
		}

		ctx := context.WithValue(context.Background(), userClaimsKey, claims)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := apiInst.RequireRole("admin")
		middleware(handler).ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("unauthorized without claims", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		w := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := apiInst.RequireRole("admin")
		middleware(handler).ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

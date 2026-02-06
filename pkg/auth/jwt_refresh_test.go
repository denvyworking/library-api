package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateRefreshToken(t *testing.T) {
	service := NewJWTService("test-secret-key")

	t.Run("success", func(t *testing.T) {
		token, err := service.GenerateRefreshToken()
		require.NoError(t, err)
		require.NotEmpty(t, token)
		require.Greater(t, len(token), 20) // токен должен быть достаточно длинным
	})

	t.Run("tokens are unique", func(t *testing.T) {
		token1, err1 := service.GenerateRefreshToken()
		require.NoError(t, err1)

		token2, err2 := service.GenerateRefreshToken()
		require.NoError(t, err2)

		require.NotEqual(t, token1, token2) // токены должны быть разными
	})
}

func TestGetRefreshTokenTTL(t *testing.T) {
	service := NewJWTService("test-secret-key")

	ttl := service.GetRefreshTokenTTL()
	expectedTTL := 7 * 24 * time.Hour

	require.Equal(t, expectedTTL, ttl)
	require.Equal(t, 168*time.Hour, ttl) // 7 дней = 168 часов
}

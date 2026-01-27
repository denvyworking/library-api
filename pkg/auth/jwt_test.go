package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTService(t *testing.T) {
	service := NewJWTService("test-secret-key")

	token, err := service.GenerateAccessToken(123, "admin")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := service.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, 123, claims.UserID)
	require.Equal(t, "admin", claims.Role)

	require.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

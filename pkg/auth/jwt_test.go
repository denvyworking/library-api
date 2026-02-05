package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

	// invalid token
	_, err = service.ParseToken("invalid.token.here")
	require.Error(t, err)

	otherService := NewJWTService("wrong-secret")
	token, _ = otherService.GenerateAccessToken(123, "admin")
	_, err = service.ParseToken(token)
	require.Error(t, err)
}

func TestJWTService_ExpiredToken(t *testing.T) {
	claims := Claims{
		UserID: 123,
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("test-secret-key"))

	service := NewJWTService("test-secret-key")
	_, err := service.ParseToken(tokenStr)
	require.Error(t, err)
}

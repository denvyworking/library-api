package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hashed, err := HashPassword("mypassword")
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = CheckPassword(hashed, "mypassword")
	require.NoError(t, err)

	err = CheckPassword(hashed, "wrongpassword")
	require.Error(t, err)
}

func TestHashPassword_Empty(t *testing.T) {
	hashed, err := HashPassword("")
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = CheckPassword(hashed, "")
	require.NoError(t, err)
}

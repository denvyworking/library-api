package service

import (
	"context"
	"leti/pkg/auth"
	"leti/pkg/models"
	"leti/pkg/repository/fake"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	fakeDB := &fake.FakeRepo{}
	srv := NewService(fakeDB)

	t.Run("success with default role", func(t *testing.T) {
		userID, err := srv.RegisterUser(context.Background(), "newuser", "password123", "")
		require.NoError(t, err)
		require.Greater(t, userID, 0)

		user, err := fakeDB.GetUserByUsername(context.Background(), "newuser")
		require.NoError(t, err)
		require.Equal(t, "newuser", user.Username)
		require.Equal(t, "user", user.Role)
		require.NotEqual(t, "password123", user.Password)
	})

	t.Run("success with admin role", func(t *testing.T) {
		userID, err := srv.RegisterUser(context.Background(), "adminuser", "password123", "admin")
		require.NoError(t, err)
		require.Greater(t, userID, 0)

		user, err := fakeDB.GetUserByUsername(context.Background(), "adminuser")
		require.NoError(t, err)
		require.Equal(t, "admin", user.Role)
	})

	t.Run("username already exists", func(t *testing.T) {
		hashedPass, _ := auth.HashPassword("password")
		fakeDB.Users = []models.User{
			{ID: 1, Username: "existing", Password: hashedPass, Role: "user"},
		}

		_, err := srv.RegisterUser(context.Background(), "existing", "password123", "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})

	t.Run("invalid role defaults to user", func(t *testing.T) {
		_, err := srv.RegisterUser(context.Background(), "invalidrole", "password123", "invalid")
		require.NoError(t, err)

		user, err := fakeDB.GetUserByUsername(context.Background(), "invalidrole")
		require.NoError(t, err)
		require.Equal(t, "user", user.Role)
	})
}

func TestRefreshTokens(t *testing.T) {
	fakeDB := &fake.FakeRepo{}
	srv := NewService(fakeDB)

	// Создаем пользователя
	hashedPass, _ := auth.HashPassword("password")
	fakeDB.Users = []models.User{
		{ID: 1, Username: "testuser", Password: hashedPass, Role: "user"},
	}

	t.Run("success", func(t *testing.T) {
		expiresAt := time.Now().Add(24 * time.Hour)
		err := fakeDB.SaveRefreshToken(context.Background(), 1, "valid-refresh-token", expiresAt)
		require.NoError(t, err)

		user, err := srv.RefreshTokens(context.Background(), "valid-refresh-token")
		require.NoError(t, err)
		require.Equal(t, 1, user.ID)
		require.Equal(t, "testuser", user.Username)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := srv.RefreshTokens(context.Background(), "invalid-token")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid or expired")
	})

	t.Run("expired token", func(t *testing.T) {
		// Сохраняем истекший токен
		expiresAt := time.Now().Add(-1 * time.Hour)
		fakeDB.SaveRefreshToken(context.Background(), 1, "expired-token", expiresAt)

		_, err := srv.RefreshTokens(context.Background(), "expired-token")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid or expired")
	})

	t.Run("user not found", func(t *testing.T) {
		// Токен для несуществующего пользователя
		expiresAt := time.Now().Add(24 * time.Hour)
		fakeDB.SaveRefreshToken(context.Background(), 999, "token-for-missing-user", expiresAt)

		_, err := srv.RefreshTokens(context.Background(), "token-for-missing-user")
		require.Error(t, err)
		require.Contains(t, err.Error(), "user not found")
	})
}

func TestSaveRefreshToken(t *testing.T) {
	fakeDB := &fake.FakeRepo{}
	srv := NewService(fakeDB)

	t.Run("success", func(t *testing.T) {
		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		err := srv.SaveRefreshToken(context.Background(), 1, "test-token", expiresAt)
		require.NoError(t, err)
		rt, err := fakeDB.GetRefreshToken(context.Background(), "test-token")
		require.NoError(t, err)
		require.Equal(t, 1, rt.UserID)
		require.Equal(t, "test-token", rt.Token)
	})
}

func TestDeleteRefreshToken(t *testing.T) {
	fakeDB := &fake.FakeRepo{}
	srv := NewService(fakeDB)

	t.Run("success", func(t *testing.T) {
		expiresAt := time.Now().Add(24 * time.Hour)
		fakeDB.SaveRefreshToken(context.Background(), 1, "token-to-delete", expiresAt)
		err := srv.DeleteRefreshToken(context.Background(), "token-to-delete")
		require.NoError(t, err)
		_, err = fakeDB.GetRefreshToken(context.Background(), "token-to-delete")
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("delete non-existent token", func(t *testing.T) {
		err := srv.DeleteRefreshToken(context.Background(), "non-existent")
		require.NoError(t, err)
	})
}

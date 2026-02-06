package fake

import (
	"context"
	"leti/pkg/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFakeRepo_UserDB(t *testing.T) {
	repo := &FakeRepo{
		Users: []models.User{
			{ID: 1, Username: "user1", Password: "hash1", Role: "user"},
			{ID: 2, Username: "admin1", Password: "hash2", Role: "admin"},
		},
	}

	t.Run("GetUserByUsername success", func(t *testing.T) {
		user, err := repo.GetUserByUsername(context.Background(), "user1")
		require.NoError(t, err)
		require.Equal(t, 1, user.ID)
		require.Equal(t, "user1", user.Username)
	})

	t.Run("GetUserByUsername not found", func(t *testing.T) {
		_, err := repo.GetUserByUsername(context.Background(), "nonexistent")
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("GetUserByID success", func(t *testing.T) {
		user, err := repo.GetUserByID(context.Background(), 2)
		require.NoError(t, err)
		require.Equal(t, 2, user.ID)
		require.Equal(t, "admin1", user.Username)
	})

	t.Run("GetUserByID not found", func(t *testing.T) {
		_, err := repo.GetUserByID(context.Background(), 999)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("CreateUser success", func(t *testing.T) {
		newUser := models.User{
			Username: "newuser",
			Password: "hash3",
			Role:     "user",
		}

		id, err := repo.CreateUser(context.Background(), newUser)
		require.NoError(t, err)
		require.Equal(t, 3, id)

		user, err := repo.GetUserByID(context.Background(), 3)
		require.NoError(t, err)
		require.Equal(t, "newuser", user.Username)
	})

	t.Run("CreateUser duplicate username", func(t *testing.T) {
		newUser := models.User{
			Username: "user1",
			Password: "hash4",
			Role:     "user",
		}

		_, err := repo.CreateUser(context.Background(), newUser)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
	})
}

func TestFakeRepo_RefreshTokens(t *testing.T) {
	repo := &FakeRepo{
		Users: []models.User{
			{ID: 1, Username: "user1", Password: "hash1", Role: "user"},
		},
	}

	t.Run("SaveRefreshToken success", func(t *testing.T) {
		expiresAt := time.Now().Add(24 * time.Hour)
		err := repo.SaveRefreshToken(context.Background(), 1, "token1", expiresAt)
		require.NoError(t, err)

		rt, err := repo.GetRefreshToken(context.Background(), "token1")
		require.NoError(t, err)
		require.Equal(t, 1, rt.UserID)
		require.Equal(t, "token1", rt.Token)
	})

	t.Run("GetRefreshToken success", func(t *testing.T) {
		expiresAt := time.Now().Add(24 * time.Hour)
		repo.SaveRefreshToken(context.Background(), 1, "valid-token", expiresAt)

		rt, err := repo.GetRefreshToken(context.Background(), "valid-token")
		require.NoError(t, err)
		require.Equal(t, 1, rt.UserID)
		require.Equal(t, "valid-token", rt.Token)
	})

	t.Run("GetRefreshToken expired", func(t *testing.T) {
		expiresAt := time.Now().Add(-1 * time.Hour) // истекший токен
		repo.SaveRefreshToken(context.Background(), 1, "expired-token", expiresAt)

		_, err := repo.GetRefreshToken(context.Background(), "expired-token")
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found or expired")
	})

	t.Run("GetRefreshToken not found", func(t *testing.T) {
		_, err := repo.GetRefreshToken(context.Background(), "nonexistent")
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found or expired")
	})

	t.Run("DeleteRefreshToken success", func(t *testing.T) {
		expiresAt := time.Now().Add(24 * time.Hour)
		repo.SaveRefreshToken(context.Background(), 1, "token-to-delete", expiresAt)

		err := repo.DeleteRefreshToken(context.Background(), "token-to-delete")
		require.NoError(t, err)
		_, err = repo.GetRefreshToken(context.Background(), "token-to-delete")
		require.Error(t, err)
	})

	t.Run("DeleteUserRefreshTokens success", func(t *testing.T) {
		expiresAt := time.Now().Add(24 * time.Hour)
		repo.SaveRefreshToken(context.Background(), 1, "token1", expiresAt)
		repo.SaveRefreshToken(context.Background(), 1, "token2", expiresAt)
		repo.SaveRefreshToken(context.Background(), 2, "token3", expiresAt) // другой пользователь

		err := repo.DeleteUserRefreshTokens(context.Background(), 1)
		require.NoError(t, err)

		_, err = repo.GetRefreshToken(context.Background(), "token1")
		require.Error(t, err)

		_, err = repo.GetRefreshToken(context.Background(), "token2")
		require.Error(t, err)

		rt, err := repo.GetRefreshToken(context.Background(), "token3")
		require.NoError(t, err)
		require.Equal(t, 2, rt.UserID)
	})
}

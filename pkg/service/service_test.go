package service

import (
	"context"
	"leti/pkg/auth"
	"leti/pkg/models"
	"leti/pkg/repository/fake"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateUserCredentials(t *testing.T) {
	// Создаём фейковую БД
	fakeDB := &fake.FakeRepo{}

	// Добавляем пользователя напрямую в фейковую БД
	hashedPass, _ := auth.HashPassword("password")
	fakeDB.Users = []models.User{
		{ID: 1, Username: "Den", Password: hashedPass, Role: "admin"},
	}

	srv := NewService(fakeDB)

	t.Run("success", func(t *testing.T) {
		user, err := srv.ValidateUserCredentials(context.Background(), "Den", "password")
		require.NoError(t, err)
		require.Equal(t, "admin", user.Role)
	})

	t.Run("wrong password", func(t *testing.T) {
		_, err := srv.ValidateUserCredentials(context.Background(), "Den", "wrong")
		require.Error(t, err)
	})
}

func TestNewAuthor(t *testing.T) {
	fakeDB := &fake.FakeRepo{}
	srv := NewService(fakeDB)

	t.Run("empty name", func(t *testing.T) {
		_, err := srv.NewAuthor(context.Background(), models.Author{Author: ""})
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("valid author", func(t *testing.T) {
		id, err := srv.NewAuthor(context.Background(), models.Author{Author: "Толстой"})
		require.NoError(t, err)
		require.Equal(t, 1, id)

		// Проверяем, что автор действительно создан
		authors, _ := fakeDB.GetAllAuthors(context.Background())
		require.Len(t, authors, 1)
		require.Equal(t, "Толстой", authors[0].Author)
	})
}

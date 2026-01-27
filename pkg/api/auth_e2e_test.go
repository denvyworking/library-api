package api

import (
	"leti/pkg/api/dto"
	"leti/pkg/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_JWTAuthorization(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Создаём тестового пользователя (если нужно)
	// Но у нас уже есть Den из миграций

	t.Run("successful login", func(t *testing.T) {
		token := login(t, ts.URL, "Den", "password")
		require.NotEmpty(t, token)
	})

	t.Run("debug - check users", func(t *testing.T) {
		token := login(t, ts.URL, "Den", "password")
		require.NotEmpty(t, token)
	})
	t.Run("create book with valid token", func(t *testing.T) {
		token := login(t, ts.URL, "Den", "password")
		authorID := createAuthor(t, ts.URL, "Толстой_unique")
		genreID := createGenre(t, ts.URL, "Роман_unique")

		bookID := createBook(t, ts.URL, token, dto.CreateBookRequest{
			Name:     "Война и мир",
			AuthorID: authorID,
			GenreID:  genreID,
			Price:    1000,
		})
		require.Greater(t, bookID, 0)
	})

	t.Run("create book without token", func(t *testing.T) {
		authorID := createAuthor(t, ts.URL, "Достоевский")
		genreID := createGenre(t, ts.URL, "Роман")

		body := marshal(t, dto.CreateBookRequest{
			Name:     "Преступление и наказание",
			AuthorID: authorID,
			GenreID:  genreID,
			Price:    800,
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/books", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("create book with invalid token", func(t *testing.T) {
		authorID := createAuthor(t, ts.URL, "Чехов")
		genreID := createGenre(t, ts.URL, "Рассказ")

		body := marshal(t, dto.CreateBookRequest{
			Name:     "Вишнёвый сад",
			AuthorID: authorID,
			GenreID:  genreID,
			Price:    600,
		})

		req := newRequestWithAuth(t, http.MethodPost, ts.URL+"/api/books", "invalid-token", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

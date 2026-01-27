////go:build e2e

package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"leti/pkg/api/dto"
	"leti/pkg/models"
	"leti/pkg/service"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
)

func TestE2E_FullBookFlow(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()
	token := login(t, ts.URL, "Den", "password")

	testCases := []struct {
		booksNumber int
		author      string
		genre       string
		name        string
		price       int
	}{
		{1, "Антон Чехов", "Рассказ", "Вишнёвый сад", 750},
		{2, "Александр Пушкин", "Роман", "Евгений Онегин", 400},
	}

	for _, tc := range testCases {
		t.Run(strconv.Itoa(tc.booksNumber), func(t *testing.T) {
			authorID := createAuthor(t, ts.URL, tc.author)
			genreID := createGenre(t, ts.URL, tc.genre)

			bookID := createBook(t, ts.URL, token, dto.CreateBookRequest{
				Name:     tc.name,
				AuthorID: authorID,
				GenreID:  genreID,
				Price:    tc.price,
			})

			books := getBooksWithAuthors(t, ts.URL, token)

			require.Len(t, books, 1)
			require.Equal(t, tc.name, books[0].Name)
			require.Equal(t, tc.author, books[0].AuthorName)
			require.Equal(t, tc.price, books[0].Price)

			deleteBook(t, ts.URL, token, bookID)

			books = getBooksWithAuthors(t, ts.URL, token)
			require.Len(t, books, 0)

		})
	}

}

func TestE2E_InvalidJSON(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()
	token := login(t, ts.URL, "Den", "password")

	body := []byte(`{"name": "Book", "author_id": "not-a-number", "genre_id": 1, "price": 100}`)
	req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestE2E_DeleteBookByWrongID(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()
	token := login(t, ts.URL, "Den", "password")

	authorID := createAuthor(t, ts.URL, "Иван Тургеньев")
	genreID := createGenre(t, ts.URL, "Рассказ")

	createBook(t, ts.URL, token, dto.CreateBookRequest{
		Name:     "МуМу",
		AuthorID: authorID,
		GenreID:  genreID,
		Price:    100,
	})

	deleteBook(t, ts.URL, token, 999)

	books := getBooks(t, ts.URL, token)
	require.Len(t, books, 1)
}

// === ТЕСТ: авторизация ===
func TestE2E_Authorization(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()
	token := login(t, ts.URL, "Den", "password")

	// 1. Без токена -> Get 200 OK
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/books", nil)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 2. С неверным токеном -> 401
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+"/api/books?id=999", nil)
	req.Header.Set("Authorization", "wrong-token")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

	// 3. С верным токеном -> 200
	req, _ = http.NewRequest(http.MethodGet, ts.URL+"/api/books", nil)
	req.Header.Set("Authorization", token)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 4. Не метод Get без токена -> 401
	req, _ = http.NewRequest(http.MethodPost, ts.URL+"/api/books", nil)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()

}

func TestE2E_GetBookById(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()
	token := login(t, ts.URL, "Den", "password")

	testCases := []struct {
		booksNumber int
		author      string
		genre       string
		name        string
		price       int
	}{
		{1, "Антон Чехов", "Рассказ", "Вишнёвый сад", 750},
		{2, "Александр Пушкин", "Роман", "Евгений Онегин", 400},
	}

	for _, tc := range testCases {
		t.Run(strconv.Itoa(tc.booksNumber), func(t *testing.T) {
			authorID := createAuthor(t, ts.URL, tc.author)
			genreID := createGenre(t, ts.URL, tc.genre)

			bookID := createBook(t, ts.URL, token, dto.CreateBookRequest{
				Name:     tc.name,
				AuthorID: authorID,
				GenreID:  genreID,
				Price:    tc.price,
			})

			book := getBookById(t, ts.URL, strconv.Itoa(bookID))
			require.Equal(t, tc.name, book.Name)
			require.Equal(t, tc.price, book.Price)

		})
	}
}

func TestE2E_PatchBook(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()
	token := login(t, ts.URL, "Den", "password")

	authorID := createAuthor(t, ts.URL, "Достоевский")
	genreID := createGenre(t, ts.URL, "Роман")

	bookID := createBook(t, ts.URL, token, dto.CreateBookRequest{
		Name:     "Идиот",
		AuthorID: authorID,
		GenreID:  genreID,
		Price:    599,
	})

	// Обновляем только цену
	patchBody := marshal(t, models.BookUpdate{
		Price: ptr(599),
	})
	req := newRequestWithAuth(t, http.MethodPatch, ts.URL+"/api/books?id="+strconv.Itoa(bookID), token, patchBody)
	resp := doRequest(t, req)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Проверяем
	books := getBooks(t, ts.URL, token)
	require.Len(t, books, 1)
	require.Equal(t, 599, books[0].Price)
	require.Equal(t, "Идиот", books[0].Name) // имя не изменилось!
}

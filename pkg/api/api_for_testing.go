package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"leti/pkg/models"
	psg "leti/pkg/repository/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"
)

func getTestDBConn() string {
	if conn := os.Getenv("TEST_DATABASE_URL"); conn != "" {
		return conn
	}
	return "postgres://postgres@localhost:5433/leti_test?sslmode=disable"
}

func setupTestDBWithMigrations(t *testing.T) *psg.PGRepo {
	t.Helper()

	connStr := getTestDBConn()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsURL := "file://" + filepath.ToSlash(filepath.Join(dir, "..", "..", "migrations"))
	m, err := migrate.New(migrationsURL, connStr)
	require.NoError(t, err)
	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	repo, err := psg.New(connStr)
	require.NoError(t, err)

	err = repo.TruncateAll(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		repo.Close()
	})

	return repo
}

func createAuthor(t *testing.T, baseURL, name string) int {
	body := marshal(t, models.Author{Author: name})
	req := newRequest(t, http.MethodPost, baseURL+"/api/authors", body)
	resp := doRequest(t, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var id int
	json.NewDecoder(resp.Body).Decode(&id)
	return id
}

func createGenre(t *testing.T, baseURL, name string) int {
	body := marshal(t, models.Genre{Genre: name})
	req := newRequest(t, http.MethodPost, baseURL+"/api/genres", body)
	resp := doRequest(t, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var id int
	json.NewDecoder(resp.Body).Decode(&id)
	return id
}

func createBook(t *testing.T, baseURL, token string, book models.Book) int {
	body := marshal(t, book)
	req := newRequestWithAuth(t, http.MethodPost, baseURL+"/api/books", token, body)
	resp := doRequest(t, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var id int
	json.NewDecoder(resp.Body).Decode(&id)
	return id
}

func deleteBook(t *testing.T, baseURL, token string, id int) {
	req := newRequestWithAuth(t, http.MethodDelete, baseURL+"/api/books?id="+strconv.Itoa(id), token, nil)
	resp := doRequest(t, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func getBooksWithAuthors(t *testing.T, baseURL, token string) []models.BookWithAuthor {
	req := newRequestWithAuth(t, http.MethodGet, baseURL+"/api/books/withauthors", token, nil)
	resp := doRequest(t, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var books []models.BookWithAuthor
	json.NewDecoder(resp.Body).Decode(&books)
	return books
}

func getBooks(t *testing.T, baseURL, token string) []models.Book {
	req := newRequestWithAuth(t, http.MethodGet, baseURL+"/api/books", token, nil)
	resp := doRequest(t, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var books []models.Book
	json.NewDecoder(resp.Body).Decode(&books)
	return books
}

func getAuthors(t *testing.T, baseURL string) []models.Author {
	req := newRequest(t, http.MethodGet, baseURL+"/api/authors", nil)
	resp := doRequest(t, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var authors []models.Author
	json.NewDecoder(resp.Body).Decode(&authors)
	return authors
}

func getGenres(t *testing.T, baseURL string) []models.Genre {
	req := newRequest(t, http.MethodGet, baseURL+"/api/genres", nil)
	resp := doRequest(t, req)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var genres []models.Genre
	json.NewDecoder(resp.Body).Decode(&genres)
	return genres
}

func newRequestWithAuth(t *testing.T, method, url, token string, body []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	require.NoError(t, err)
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func newRequest(t *testing.T, method, url string, body []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func doRequest(t *testing.T, req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func marshal(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return data
}

func ptr[T any](v T) *T { return &v }

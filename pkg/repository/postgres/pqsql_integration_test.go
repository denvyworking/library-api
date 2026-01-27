////go:build integration

package postgres

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"leti/pkg/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"
)

func getTestDBConn() string {
	if conn := os.Getenv("TEST_DATABASE_URL"); conn != "" {
		return conn
	}
	return "postgres://postgres:postgres@localhost:45432/leti_test?sslmode=disable"
}

// setupTestDB применяет миграции и возвращает репозиторий для тестов
func setupTestDB(t *testing.T) *PGRepo {
	t.Helper()

	connStr := getTestDBConn()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsURL := "file://" + filepath.ToSlash(filepath.Join(dir, "..", "..", "..", "migrations"))
	m, err := migrate.New(migrationsURL, connStr)
	require.NoError(t, err)
	require.NoError(t, m.Down())
	require.NoError(t, m.Up())

	repo, err := New(connStr)
	require.NoError(t, err)

	err = repo.TruncateAll(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		repo.Close()
	})

	return repo
}

// === ТЕСТЫ ===

func TestPGRepo_NewAuthorAndGetAll(t *testing.T) {
	repo := setupTestDB(t)

	id, err := repo.NewAuthor(context.Background(), models.Author{Author: "Лев Толстой"})
	require.NoError(t, err)
	require.Equal(t, 1, id)

	authors, err := repo.GetAllAuthors(context.Background())
	require.NoError(t, err)
	require.Len(t, authors, 1)
	require.Equal(t, "Лев Толстой", authors[0].Author)
	require.Equal(t, int(1), authors[0].ID)
}

func TestPGRepo_NewGenreAndGetAll(t *testing.T) {
	repo := setupTestDB(t)

	id, err := repo.NewGenre(context.Background(), models.Genre{Genre: "Новелла"})
	require.NoError(t, err)
	require.Equal(t, 1, id)

	genres, err := repo.GetAllGenres(context.Background())
	require.NoError(t, err)
	require.Len(t, genres, 1)
	require.Equal(t, "Новелла", genres[0].Genre)
	require.Equal(t, int(1), genres[0].ID)
}

func TestPGRepo_NewGenre_Duplicate(t *testing.T) {
	repo := setupTestDB(t)

	id1, err := repo.NewGenre(context.Background(), models.Genre{Genre: "Роман"})
	require.NoError(t, err)
	require.Equal(t, 1, id1)

	_, err = repo.NewGenre(context.Background(), models.Genre{Genre: "Роман"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unique") // или "duplicate key"
}

// могут быть авторы с одной фамилией!
func TestPGRepo_NewAuthor_NoDuplicate(t *testing.T) {
	repo := setupTestDB(t)

	id1, err := repo.NewAuthor(context.Background(), models.Author{Author: "Гоголь"})
	require.NoError(t, err)
	require.Equal(t, 1, id1)

	_, err = repo.NewAuthor(context.Background(), models.Author{Author: "Гоголь"})
	require.NoError(t, err)
}

func TestPGRepo_BooksWithAuthors(t *testing.T) {
	repo := setupTestDB(t)

	authorID, err := repo.NewAuthor(context.Background(), models.Author{Author: "Фёдор Достоевский"})
	require.NoError(t, err)

	genreID, err := repo.NewGenre(context.Background(), models.Genre{Genre: "Роман"})
	require.NoError(t, err)

	bookID, err := repo.NewBook(context.Background(), models.Book{
		Name:      "Преступление и наказание",
		Author_id: authorID,
		Genre_id:  genreID,
		Price:     999,
	})
	require.NoError(t, err)
	require.Equal(t, 1, bookID)

	books, err := repo.GetAllWithAuthors(context.Background())
	require.NoError(t, err)
	require.Len(t, books, 1)
	require.Equal(t, "Преступление и наказание", books[0].Name)
	require.Equal(t, "Фёдор Достоевский", books[0].AuthorName)
	require.Equal(t, 999, books[0].Price)
}

func TestPGRepo_BookCRUD(t *testing.T) {
	repo := setupTestDB(t)

	authorID, _ := repo.NewAuthor(context.Background(), models.Author{Author: "Пушкин"})
	genreID, _ := repo.NewGenre(context.Background(), models.Genre{Genre: "Поэма"})

	bookID, err := repo.NewBook(context.Background(), models.Book{
		Name:      "Евгений Онегин",
		Author_id: authorID,
		Genre_id:  genreID,
		Price:     500,
	})
	require.NoError(t, err)

	book, err := repo.GetBookByID(context.Background(), bookID)
	require.NoError(t, err)
	require.Equal(t, "Евгений Онегин", book.Name)
	require.Equal(t, 500, book.Price)

	newPrice := 132

	err = repo.UpdateBook(context.Background(), bookID, models.BookUpdate{
		Price: &newPrice,
	})
	require.NoError(t, err)

	book, _ = repo.GetBookByID(context.Background(), bookID)
	require.Equal(t, 132, book.Price)

	err = repo.DeleteBookById(context.Background(), bookID)
	require.NoError(t, err)

	_, err = repo.GetBookByID(context.Background(), bookID)
	require.Error(t, err) // должно быть "no rows"
}

func TestPGRepo_Book_NegativePrice(t *testing.T) {
	repo := setupTestDB(t)
	_, err := repo.NewBook(context.Background(), models.Book{
		Name: "Book", Author_id: 1, Genre_id: 1, Price: -100,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "check constraint")
}

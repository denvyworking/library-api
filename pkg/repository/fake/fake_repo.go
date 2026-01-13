// Package fake содержит фейковую реализацию repository.DataBase для unit-тестов.
package fake

import (
	"context"
	"errors"
	"fmt"
	"leti/pkg/models"
	"sync"
)

// FakeRepo реализует интерфейс repository.DataBase
type FakeRepo struct {
	mu sync.RWMutex

	// Хранилища данных (имитируют БД)
	authors []models.Author
	books   []models.Book
	genres  []models.Genre

	// Флаги для эмуляции ошибок (опционально)
	NewAuthorErr error
	NewBookErr   error
	NewGenreErr  error
}

// --- AuthorDB ---

func (f *FakeRepo) GetAllAuthors(ctx context.Context) ([]models.Author, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	// Копируем, чтобы избежать гонок при модификации
	authors := make([]models.Author, len(f.authors))
	copy(authors, f.authors)
	return authors, nil
}

func (f *FakeRepo) NewAuthor(ctx context.Context, author models.Author) (int, error) {
	if f.NewAuthorErr != nil {
		return 0, f.NewAuthorErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	id := len(f.authors) + 1
	newAuthor := models.Author{
		ID:     id,
		Author: author.Author,
	}
	f.authors = append(f.authors, newAuthor)
	return id, nil
}

// --- GenreDB ---

func (f *FakeRepo) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	genres := make([]models.Genre, len(f.genres))
	copy(genres, f.genres)
	return genres, nil
}

func (f *FakeRepo) NewGenre(ctx context.Context, genre models.Genre) (int, error) {
	if f.NewGenreErr != nil {
		return 0, f.NewGenreErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	id := len(f.genres) + 1
	newGenre := models.Genre{
		ID:    id,
		Genre: genre.Genre,
	}
	f.genres = append(f.genres, newGenre)
	return id, nil
}

// --- BooksDB ---

func (f *FakeRepo) GetBooks(ctx context.Context) ([]models.Book, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	books := make([]models.Book, len(f.books))
	copy(books, f.books)
	return books, nil
}

func (f *FakeRepo) NewBook(ctx context.Context, book models.Book) (int, error) {
	if f.NewBookErr != nil {
		return 0, f.NewBookErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	id := len(f.books) + 1
	newBook := models.Book{
		ID:        id,
		Name:      book.Name,
		Author_id: book.Author_id,
		Genre_id:  book.Genre_id,
		Price:     book.Price,
	}
	f.books = append(f.books, newBook)
	return id, nil
}

func (f *FakeRepo) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, book := range f.books {
		if int(book.ID) == id {
			return book, nil
		}
	}
	return models.Book{}, errors.New("book not found")
}

func (f *FakeRepo) DeleteBookById(ctx context.Context, id int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, book := range f.books {
		if int(book.ID) == id {
			f.books = append(f.books[:i], f.books[i+1:]...)
			return nil
		}
	}
	return errors.New("book not found")
}

func (f *FakeRepo) GetAllWithAuthors(ctx context.Context) ([]models.BookWithAuthor, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var result []models.BookWithAuthor
	for _, book := range f.books {
		// Ищем автора по Author_id
		var authorName string
		for _, author := range f.authors {
			if int(author.ID) == book.Author_id {
				authorName = author.Author
				break
			}
		}
		result = append(result, models.BookWithAuthor{
			ID:         book.ID,
			Name:       book.Name,
			Price:      book.Price,
			GenreID:    -1,
			AuthorID:   book.Author_id,
			AuthorName: authorName,
		})
	}
	return result, nil
}

func (f *FakeRepo) UpdateBook(ctx context.Context, id int, update models.BookUpdate) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, book := range f.books {
		if book.ID == id {
			if update.Name != nil {
				f.books[i].Name = *update.Name
			}
			if update.Price != nil {
				if *update.Price < 0 {
					return errors.New("price must be non-negative")
				}
				f.books[i].Price = *update.Price
			}
			return nil
		}
	}
	return fmt.Errorf("book with id %d not found", id)
}

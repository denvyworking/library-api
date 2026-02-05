// Package fake содержит фейковую реализацию repository.DataBase для unit-тестов.
package fake

import (
	"context"
	"errors"
	"fmt"
	"leti/pkg/models"
	"sync"
	"time"
)

// FakeRepo реализует интерфейс repository.DataBase
type FakeRepo struct {
	mu sync.RWMutex

	// Хранилища данных (имитируют БД)
	Authors      []models.Author
	Books        []models.Book
	Genres       []models.Genre
	Users        []models.User
	RefreshToken []models.RefreshToken

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
	authors := make([]models.Author, len(f.Authors))
	copy(authors, f.Authors)
	return authors, nil
}

func (f *FakeRepo) NewAuthor(ctx context.Context, author models.Author) (int, error) {
	if f.NewAuthorErr != nil {
		return 0, f.NewAuthorErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	id := len(f.Authors) + 1
	newAuthor := models.Author{
		ID:     id,
		Author: author.Author,
	}
	f.Authors = append(f.Authors, newAuthor)
	return id, nil
}

// --- GenreDB ---

func (f *FakeRepo) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	genres := make([]models.Genre, len(f.Genres))
	copy(genres, f.Genres)
	return genres, nil
}

func (f *FakeRepo) NewGenre(ctx context.Context, genre models.Genre) (int, error) {
	if f.NewGenreErr != nil {
		return 0, f.NewGenreErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	id := len(f.Genres) + 1
	newGenre := models.Genre{
		ID:    id,
		Genre: genre.Genre,
	}
	f.Genres = append(f.Genres, newGenre)
	return id, nil
}

// --- BooksDB ---

func (f *FakeRepo) GetBooks(ctx context.Context) ([]models.Book, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	books := make([]models.Book, len(f.Books))
	copy(books, f.Books)
	return books, nil
}

func (f *FakeRepo) NewBook(ctx context.Context, book models.Book) (int, error) {
	if f.NewBookErr != nil {
		return 0, f.NewBookErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	id := len(f.Books) + 1
	newBook := models.Book{
		ID:        id,
		Name:      book.Name,
		Author_id: book.Author_id,
		Genre_id:  book.Genre_id,
		Price:     book.Price,
	}
	f.Books = append(f.Books, newBook)
	return id, nil
}

func (f *FakeRepo) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, book := range f.Books {
		if int(book.ID) == id {
			return book, nil
		}
	}
	return models.Book{}, errors.New("book not found")
}

func (f *FakeRepo) DeleteBookById(ctx context.Context, id int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, book := range f.Books {
		if int(book.ID) == id {
			f.Books = append(f.Books[:i], f.Books[i+1:]...)
			return nil
		}
	}
	return errors.New("book not found")
}

func (f *FakeRepo) GetAllWithAuthors(ctx context.Context) ([]models.BookWithAuthor, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var result []models.BookWithAuthor
	for _, book := range f.Books {
		var authorName string
		for _, author := range f.Authors {
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

	for i, book := range f.Books {
		if book.ID == id {
			if update.Name != nil {
				f.Books[i].Name = *update.Name
			}
			if update.Price != nil {
				if *update.Price < 0 {
					return errors.New("price must be non-negative")
				}
				f.Books[i].Price = *update.Price
			}
			return nil
		}
	}
	return fmt.Errorf("book with id %d not found", id)
}

// --- UserDB ---

func (f *FakeRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for i := range f.Users {
		if f.Users[i].Username == username {
			u := f.Users[i]
			return &u, nil
		}
	}
	return nil, fmt.Errorf("the user was not found")
}

func (f *FakeRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for i := range f.Users {
		if f.Users[i].ID == userID {
			u := f.Users[i]
			return &u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (f *FakeRepo) CreateUser(ctx context.Context, user models.User) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, u := range f.Users {
		if u.Username == user.Username {
			return 0, errors.New("username already exists")
		}
	}

	id := len(f.Users) + 1
	user.ID = id
	f.Users = append(f.Users, user)
	return id, nil
}

func (f *FakeRepo) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := len(f.RefreshToken) - 1; i >= 0; i-- {
		if f.RefreshToken[i].Token == token {
			f.RefreshToken = append(f.RefreshToken[:i], f.RefreshToken[i+1:]...)
			break
		}
	}

	id := len(f.RefreshToken) + 1
	f.RefreshToken = append(f.RefreshToken, models.RefreshToken{
		ID:        id,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	})
	return nil
}

func (f *FakeRepo) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	now := time.Now()
	for i := range f.RefreshToken {
		rt := f.RefreshToken[i]
		if rt.Token == token && rt.ExpiresAt.After(now) {
			copyRT := rt
			return &copyRT, nil
		}
	}
	return nil, errors.New("refresh token not found or expired")
}

func (f *FakeRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, rt := range f.RefreshToken {
		if rt.Token == token {
			f.RefreshToken = append(f.RefreshToken[:i], f.RefreshToken[i+1:]...)
			return nil
		}
	}
	return nil
}

func (f *FakeRepo) DeleteUserRefreshTokens(ctx context.Context, userID int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	var newList []models.RefreshToken
	for _, rt := range f.RefreshToken {
		if rt.UserID != userID {
			newList = append(newList, rt)
		}
	}
	f.RefreshToken = newList
	return nil
}

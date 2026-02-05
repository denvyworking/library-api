package repository

import (
	"context"
	"leti/pkg/models"
	"time"
)

type AuthorDB interface {
	GetAllAuthors(context.Context) ([]models.Author, error)
	NewAuthor(context.Context, models.Author) (int, error)
}

type BooksDB interface {
	GetBooks(context.Context) ([]models.Book, error)
	NewBook(context.Context, models.Book) (int, error)
	GetBookByID(context.Context, int) (models.Book, error)
	DeleteBookById(context.Context, int) error
	GetAllWithAuthors(context.Context) ([]models.BookWithAuthor, error)
	UpdateBook(context.Context, int, models.BookUpdate) error
}

type GenreDB interface {
	GetAllGenres(context.Context) ([]models.Genre, error)
	NewGenre(context.Context, models.Genre) (int, error)
}

type UserDB interface {
	GetUserByUsername(context.Context, string) (*models.User, error)
	GetUserByID(context.Context, int) (*models.User, error)
	CreateUser(context.Context, models.User) (int, error)
	// старый refresh убрать, новый сохранить
	SaveRefreshToken(context.Context, int, string, time.Time) error
	// проверить токен и срок
	GetRefreshToken(context.Context, string) (*models.RefreshToken, error)
	DeleteRefreshToken(context.Context, string) error
	DeleteUserRefreshTokens(context.Context, int) error
}

type DataBase interface {
	BooksDB
	GenreDB
	AuthorDB
	UserDB
}

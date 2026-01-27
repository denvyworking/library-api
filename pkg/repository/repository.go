package repository

import (
	"context"
	"leti/pkg/models"
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
}

type DataBase interface {
	BooksDB
	GenreDB
	AuthorDB
	UserDB
}

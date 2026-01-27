package service

import (
	"context"
	"errors"
	"fmt"
	"leti/pkg/models"
	"strings"
)

func (s *Service) CreateBook(ctx context.Context, book models.Book) (int, error) {
	return s.db.NewBook(ctx, book)
}

func (s *Service) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	return s.db.GetBookByID(ctx, id)
}

func (s *Service) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return s.db.GetBooks(ctx) // Передаем контекст дальше
}

func (s *Service) RemoveBook(ctx context.Context, id int) error {
	return s.db.DeleteBookById(ctx, id)
}

func (s *Service) UpdateBook(ctx context.Context, id int, update models.BookUpdate) error {
	if update.Price != nil && *update.Price < 0 {
		return errors.New("price must be non-negative")
	}
	err := s.db.UpdateBook(ctx, id, update)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("book not found: %w", err)
		}
		return err
	}
	return nil
}

func (s *Service) GetAllWithAuthors(ctx context.Context) ([]models.BookWithAuthor, error) {
	return s.db.GetAllWithAuthors(ctx)
}

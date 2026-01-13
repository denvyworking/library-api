package service

import (
	"context"
	"errors"
	"leti/pkg/models"
	"strings"
)

func (s *Service) GetAllAuthors(ctx context.Context) ([]models.Author, error) {
	return s.db.GetAllAuthors(ctx)
}

func (s *Service) NewAuthor(ctx context.Context, author models.Author) (int, error) {
	if strings.TrimSpace(author.Author) == "" {
		return 0, errors.New("author name cannot be empty")
	}
	return s.db.NewAuthor(ctx, author)
}

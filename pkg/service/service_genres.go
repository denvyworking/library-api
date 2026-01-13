package service

import (
	"context"
	"errors"
	"leti/pkg/models"
	"strings"
)

func (s *Service) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	return s.db.GetAllGenres(ctx)
}

func (s *Service) NewGenre(ctx context.Context, genre models.Genre) (int, error) {
	if strings.TrimSpace(genre.Genre) == "" {
		return 0, errors.New("genre name cannot be empty")
	}
	return s.db.NewGenre(ctx, genre)
}

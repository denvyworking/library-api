package service

import (
	"context"
	"leti/pkg/auth"
	"leti/pkg/models"
)

func (s *Service) ValidateUserCredentials(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := auth.CheckPassword(user.Password, password); err != nil {
		return nil, err
	}
	return user, nil

}

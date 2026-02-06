package service

import (
	"context"
	"errors"
	"leti/pkg/auth"
	"leti/pkg/models"
	"time"
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

// RegisterUser создает нового пользователя и возвращает его ID
func (s *Service) RegisterUser(ctx context.Context, username, password, role string) (int, error) {
	// Проверяем, существует ли пользователь
	existingUser, err := s.db.GetUserByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return 0, errors.New("username already exists")
	}

	// Хешируем пароль
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	// Устанавливаем роль по умолчанию, если не указана
	if role == "" {
		role = "user"
	}

	// Проверяем валидность роли
	if role != "user" && role != "admin" {
		role = "user"
	}

	user := models.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	id, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RefreshTokens проверяет refresh токен и возвращает пользователя для выдачи новой пары токенов
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*models.User, error) {
	// Получаем refresh токен из БД
	rt, err := s.db.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Получаем пользователя по ID из refresh токена
	user, err := s.db.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// SaveRefreshToken сохраняет refresh токен в БД
func (s *Service) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	return s.db.SaveRefreshToken(ctx, userID, token, expiresAt)
}

// DeleteRefreshToken удаляет refresh токен из БД
func (s *Service) DeleteRefreshToken(ctx context.Context, token string) error {
	return s.db.DeleteRefreshToken(ctx, token)
}

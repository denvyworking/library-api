package postgres

import (
	"context"
	"errors"
	"leti/pkg/models"
	"time"

	"github.com/jackc/pgx/v4"
)

func (repo *PGRepo) GetUserByUsername(ctx context.Context, userName string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	var user models.User
	err := repo.pool.QueryRow(ctx, `
        SELECT id, username, password, role 
        FROM users 
        WHERE username = $1;
        `,
		userName,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PGRepo) CreateUser(ctx context.Context, user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	var id int
	err := repo.pool.QueryRow(ctx, `
		INSERT INTO users (username, password, role)
		VALUES ($1, $2, $3)
		RETURNING id;
	`, user.Username, user.Password, user.Role).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *PGRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	var user models.User
	err := repo.pool.QueryRow(ctx, `
		SELECT id, username, password, role
		FROM users
		WHERE id = $1;
	`, userID).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PGRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	_, err := repo.pool.Exec(ctx, `
        DELETE FROM refresh_tokens
        WHERE token = $1;
    `, token)
	return err
}

func (repo *PGRepo) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	_, err := repo.pool.Exec(ctx, `
        INSERT INTO refresh_tokens (user_id, token, expires_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (token) DO UPDATE
        SET expires_at = EXCLUDED.expires_at;
    `, userID, token, expiresAt)
	return err
}

func (repo *PGRepo) DeleteUserRefreshTokens(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	_, err := repo.pool.Exec(ctx, `
        DELETE FROM refresh_tokens
        WHERE user_id = $1;
    `, userID)
	return err
}

func (repo *PGRepo) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	var rt models.RefreshToken
	err := repo.pool.QueryRow(ctx, `
        SELECT id, user_id, token, expires_at, created_at
        FROM refresh_tokens
        WHERE token = $1 AND expires_at > NOW();
    `, token).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("refresh token not found or expired")
		}
		return nil, err
	}
	return &rt, nil
}

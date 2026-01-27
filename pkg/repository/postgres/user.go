package postgres

import (
	"context"
	"leti/pkg/models"
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

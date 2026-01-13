package postgres

import (
	"context"

	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGRepo struct {
	pool      *pgxpool.Pool
	dbTimeout time.Duration
}

func New(connStr string) (*PGRepo, error) {
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return &PGRepo{
		pool:      pool,
		dbTimeout: 3 * time.Second}, nil
}

func (r *PGRepo) TruncateAll(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		TRUNCATE TABLE authors, genres, books RESTART IDENTITY CASCADE;
	`)
	return err
}

func (r *PGRepo) Close() {
	r.pool.Close()
}

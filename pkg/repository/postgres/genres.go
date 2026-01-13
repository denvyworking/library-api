package postgres

import (
	"context"
	"leti/pkg/models"
)

func (repo *PGRepo) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()
	rows, err := repo.pool.Query(ctx, `
		SELECT id, genre
		FROM genres
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	genres := []models.Genre{}
	for rows.Next() {
		var genre models.Genre
		err = rows.Scan(
			&genre.ID,
			&genre.Genre,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil

}

func (repo *PGRepo) NewGenre(ctx context.Context, newGenre models.Genre) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()
	var id int
	err := repo.pool.QueryRow(ctx, `
		INSERT INTO genres (genre)
		VALUES ($1)
		RETURNING id;
	`, newGenre.Genre).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

package postgres

import (
	"context"
	"leti/pkg/models"
)

func (repo *PGRepo) GetAllAuthors(ctx context.Context) ([]models.Author, error) {
	// Создаем дочерний контекст с таймаутом 3 секунды
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel() // ВСЕГДА вызываем cancel в defer!

	rows, err := repo.pool.Query(ctx, `
        SELECT id, author
        FROM authors
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	authors := []models.Author{}
	for rows.Next() {
		var author models.Author
		err = rows.Scan(
			&author.ID,
			&author.Author,
		)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	// Проверяем контекст после цикла
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (repo *PGRepo) NewAuthor(ctx context.Context, author models.Author) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()
	var id int
	err := repo.pool.QueryRow(ctx, `
		INSERT INTO authors (author)
		VALUES ($1)
		RETURNING id; 
	`, author.Author).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

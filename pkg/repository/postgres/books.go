package postgres

import (
	"context"
	"fmt"
	"leti/pkg/models"
	"strings"
)

// метод, который возвращает нам все книжки.
// Query - возвращает набор колонок
// QueryRow() - возвращает одну колонку
// Exec() - проверяет командный тег INSERT DELETE UPDATE послать через данный командый тег
func (repo *PGRepo) GetBooks(ctx context.Context) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	rows, err := repo.pool.Query(ctx, `
        SELECT id, name, author_id, genre_id, price
        FROM books
        WHERE author_id IS NOT NULL AND genre_id IS NOT NULL;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.Book
	for rows.Next() {
		var item models.Book
		if err := rows.Scan(&item.ID, &item.Name, &item.Author_id, &item.Genre_id, &item.Price); err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func (repo *PGRepo) NewBook(ctx context.Context, item models.Book) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	var id int
	// возвращает id сразу
	err := repo.pool.QueryRow(ctx, `
		INSERT INTO books (name, author_id, genre_id, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id; 
	`,
		item.Name,
		item.Author_id,
		item.Genre_id,
		item.Price,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Книга по id конкретная
func (repo *PGRepo) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()
	var book models.Book
	err := repo.pool.QueryRow(ctx, `
		SELECT id, name, author_id, genre_id, price
		FROM books
		WHERE author_id IS NOT NULL AND genre_id IS NOT NULL AND id=$1;

	`, id).Scan(
		&book.ID,
		&book.Name,
		&book.Author_id,
		&book.Genre_id,
		&book.Price,
	)

	if err != nil {
		return models.Book{}, err
	}
	return book, err
}

func (repo *PGRepo) DeleteBookById(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()
	_, err := repo.pool.Exec(ctx, `
		DELETE FROM books
		WHERE id=$1;
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGRepo) GetAllWithAuthors(ctx context.Context) ([]models.BookWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()
	rows, err := repo.pool.Query(ctx, `
        SELECT b.id, b.name, b.price, b.genre_id, b.author_id, a.author
        FROM books b
        JOIN authors a ON b.author_id = a.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.BookWithAuthor
	for rows.Next() {
		var b models.BookWithAuthor
		err := rows.Scan(&b.ID, &b.Name, &b.Price, &b.GenreID, &b.AuthorID, &b.AuthorName)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (repo *PGRepo) UpdateBook(ctx context.Context, id int, update models.BookUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, repo.dbTimeout)
	defer cancel()

	var setParts []string
	var args []interface{}
	args = append(args, id) //$1 = book ID

	argIndex := 2 //следующие параметры: $2, $3...

	if update.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *update.Name)
		argIndex++
	}

	if update.Price != nil {
		if *update.Price < 0 {
			return fmt.Errorf("price must be non-negative")
		}
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *update.Price)
	}

	if len(setParts) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE books SET %s WHERE id = $1",
		strings.Join(setParts, ", "),
	)

	result, err := repo.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("book with id %d not found", id)
	}

	return nil
}

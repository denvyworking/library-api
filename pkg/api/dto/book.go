package dto

import "leti/pkg/models"

type CreateBookRequest struct {
	Name     string `json:"name" validate:"required"`
	AuthorID int    `json:"author_id" validate:"required,min=1"`
	GenreID  int    `json:"genre_id" validate:"required,min=1"`
	Price    int    `json:"price" validate:"required,min=0"`
}

func (req CreateBookRequest) ToBookModel() models.Book {
	return models.Book{
		Name:      req.Name,
		Author_id: req.AuthorID,
		Genre_id:  req.GenreID,
		Price:     req.Price,
	}
}

type UpdateBookRequest struct {
	Name  *string `json:"name,omitempty"`
	Price *int    `json:"price,omitempty"`
}

func (req UpdateBookRequest) ToBookModel() models.BookUpdate {
	return models.BookUpdate{
		Name:  req.Name,
		Price: req.Price,
	}
}

type BookResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	AuthorID int    `json:"author_id"`
	GenreID  int    `json:"genre_id"`
	Price    int    `json:"price"`
}

type BookWithAuthorResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
	GenreID    int    `json:"genre_id"`
	Price      int    `json:"price"`
}

func FromBookModel(book models.Book) BookResponse {
	return BookResponse{
		ID:       book.ID,
		Name:     book.Name,
		AuthorID: book.Author_id,
		GenreID:  book.Genre_id,
		Price:    book.Price,
	}
}

func FromBookWithAuthorModel(book models.BookWithAuthor) BookWithAuthorResponse {
	return BookWithAuthorResponse{
		ID:         book.ID,
		Name:       book.Name,
		AuthorID:   book.AuthorID,
		AuthorName: book.AuthorName,
		GenreID:    book.GenreID,
		Price:      book.Price,
	}
}

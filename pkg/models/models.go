package models

type Book struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Author_id int    `json:"author_id"`
	Genre_id  int    `json:"genre_id"`
}

type Genre struct {
	ID    int    `json:"id"`
	Genre string `json:"genre"`
}

type Author struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
}

type BookWithAuthor struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	GenreID    int    `json:"genre_id"`
	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
}

type BookUpdate struct {
	Name  *string `json:"name,omitempty"`
	Price *int    `json:"price,omitempty"`
}

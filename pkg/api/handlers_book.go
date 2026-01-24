package api

import (
	"encoding/json"
	"leti/pkg/api/dto"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// GetBooksWithAuthors returns list of books with author names
// @Summary Получить список книг с авторами
// @Description Возвращает все книги из каталога вместе с именами авторов
// @Tags books
// @Produce json
// @Success 200 {array} dto.BookWithAuthorResponse "Список книг"
// @Router /api/books/withauthors [get]
func (api *api) booksWithAuthor(w http.ResponseWriter, r *http.Request) {
	data, err := api.srv.GetAllWithAuthors(r.Context())
	if err != nil {
		api.logger.Error("Failed to get books with authors", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response := make([]dto.BookWithAuthorResponse, len(data))
	for i, book := range data {
		response[i] = dto.FromBookWithAuthorModel(book)
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.logger.Error("Failed to encode books with authors", "error", err)
	}
}

// CreateBook creates a new book in catalog
// @Summary Создать новую книгу
// @Description Добавляет книгу в каталог (требуется авторизация)
// @Tags books
// @Accept json
// @Produce json
// @Param book body dto.CreateBookRequest true "Данные книги"
// @Success 201 {object} map[string]int "ID созданной книги"
// @Failure 400 {object} string "Невалидные данные"
// @Failure 401 {object} string "Неавторизован"
// @Router /api/books [post]
func (api *api) createBook(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "book name cannot be empty", http.StatusBadRequest)
		return
	}
	if req.AuthorID <= 0 {
		http.Error(w, "author_id must be positive", http.StatusBadRequest)
		return
	}
	if req.GenreID <= 0 {
		http.Error(w, "genre_id must be positive", http.StatusBadRequest)
		return
	}
	if req.Price < 0 {
		http.Error(w, "price must be non-negative", http.StatusBadRequest)
		return
	}

	book := req.ToBookModel()
	id, err := api.srv.CreateBook(r.Context(), book)
	if err != nil {
		api.logger.Error("Failed to create book", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		api.logger.Error("Failed to encode book ID", "error", err)
	}
}

// Update info about book
// @Summary Частично обновить информацию о книге
// @Description Обновляет указанные поля книги (требуется авторизация)
// @Tags books
// @Accept json
// @Produce json
// @Param book body dto.UpdateBookRequest true "Поля для обновления"
// @Param id query int true "ID книги"
// @Success 200 {object} dto.BookResponse "Измененная книга"
// @Failure 400 {object} string "Невалидные данные"
// @Failure 401 {object} string "Неавторизован"
func (api *api) updateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id query parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	update := req.ToBookModel()
	if err := api.srv.UpdateBook(r.Context(), id, update); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "non-negative") {
			http.Error(w, "price must be non-negative", http.StatusBadRequest)
			return
		}
		api.logger.Error("Failed to update book", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data, err := api.srv.GetBookByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		api.logger.Error("Failed to get book by ID", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response := dto.FromBookModel(data)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.logger.Error("Failed to encode book", "error", err)
	}
}

// DeleteBook deletes book by ID
// @Summary Удалить книгу из каталога
// @Description Удаляет книгу из каталога по ID (требуется авторизация)
// @Tags books
// @Param id query int true "ID книги"
// @Success 204
// @Failure 401 {object} string "Неавторизован"
// @Router /api/books [delete]
func (api *api) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := api.srv.RemoveBook(r.Context(), id); err != nil {
		api.logger.Error("Failed to delete book", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Get book by ID
// @Summary Получить книгу по ID
// @Description Получает информацию о книге по ID в каталоге
// @Tags books
// @Produce json
// @Param id query int false "ID книги"
// @Success 200 {object} dto.BookResponse
// @Failure 404 {object} string "Книга не найдена"
// @Router /api/book/{id} [get]
func (api *api) getBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	data, err := api.srv.GetBookByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		api.logger.Error("Failed to get book by ID", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response := dto.FromBookModel(data)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.logger.Error("Failed to encode book", "error", err)
	}
}

// Get all books
// @Summary Получить все книги
// @Description Получает информацию о всех книгах в каталоге
// @Tags books
// @Produce json
// @Success 200 {array} dto.BookResponse
// @Router /api/books [get]
func (api *api) getBooks(w http.ResponseWriter, r *http.Request) {
	data, err := api.srv.GetAllBooks(r.Context())
	if err != nil {
		api.logger.Error("Failed to get all books", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response := make([]dto.BookResponse, len(data))
	for i, book := range data {
		response[i] = dto.FromBookModel(book)
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		api.logger.Error("Failed to encode books", "error", err)
	}
}

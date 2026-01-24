package api

import (
	"encoding/json"
	"leti/pkg/api/dto"
	"net/http"
	"strings"
)

// Get all authors
// @Summary Получить всех авторов
// @Description Возвращает список всех авторов в каталоге
// @Tags authors
// @Produce json
// @Success 200 {array} dto.AuthorResponse
// @Router /api/authors [get]
func (api *api) getAuthors(w http.ResponseWriter, r *http.Request) {
	data, err := api.srv.GetAllAuthors(r.Context())
	if err != nil {
		api.logger.Error("Failed to get authors", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.FromAuthorModelsArray(data)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		api.logger.Error("Failed to encode authors", "error", err)
	}
}

// Create new author
// @Summary Создать нового автора
// @Description Добавляет нового автора в каталог
// @Tags authors
// @Accept json
// @Produce json
// @Param author body dto.CreateAuthorRequest true "Данные автора"
// @Success 201 {object} map[string]int "ID созданного автора"
// @Failure 400 {object} string "Невалидные данные"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/authors [post]
func (api *api) postAuthors(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAuthorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация до преобразования
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "author name cannot be empty", http.StatusBadRequest)
		return
	}

	author := req.ToAuthorModel()
	id, err := api.srv.NewAuthor(r.Context(), author)
	if err != nil {
		api.logger.Error("Failed to create author", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		api.logger.Error("Failed to encode author ID", "error", err)
	}
}

package api

import (
	"encoding/json"
	"leti/pkg/api/dto"
	"net/http"
	"strings"
)

// Get all genres
// @Summary Получить все жанры
// @Description Возвращает список всех жанров в каталоге
// @Tags genres
// @Produce json
// @Success 200 {array} dto.GenreResponse
// @Router /api/genres [get]
func (api *api) getGenres(w http.ResponseWriter, r *http.Request) {
	data, err := api.srv.GetAllGenres(r.Context())
	if err != nil {
		api.logger.Error("Failed to get genres", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := dto.FromGenreModelsArray(data)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		api.logger.Error("Failed to encode genres", "error", err)
	}
}

// Create new genre
// @Summary Создать новый жанр
// @Description Добавляет новый жанр в каталог
// @Tags genres
// @Accept json
// @Produce json
// @Param genre body dto.CreateGenreRequest true "Данные о жанре"
// @Success 201 {object} map[string]int "ID созданного жанра"
// @Failure 400 {object} string "Невалидные данные"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /api/genres [post]
func (api *api) postGenres(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "genre name cannot be empty", http.StatusBadRequest)
		return
	}

	genre := req.ToGenreModel()
	id, err := api.srv.NewGenre(r.Context(), genre)
	if err != nil {
		api.logger.Error("Failed to create genre", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		api.logger.Error("Failed to encode genre ID", "error", err)
	}
}

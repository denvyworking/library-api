package api

import (
	"encoding/json"
	"leti/pkg/auth"
	"net/http"
	"strings"
)

// Login handles user authentication
// @Summary Аутентификация пользователя
// @Description Возвращает JWT токен при успешной аутентификации
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body auth.LoginRequest true "Учётные данные"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} string "Невалидные данные"
// @Failure 401 {object} string "Неверные учётные данные"
// @Router /api/auth/login [post]
func (api *api) login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(req.Username) == "" {
		http.Error(w, "username cannot be empty", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		http.Error(w, "password cannot be empty", http.StatusBadRequest)
	}

	user, err := api.srv.ValidateUserCredentials(r.Context(), req.Username, req.Password)
	if err != nil {
		api.logger.Error("Invalid credentials", "username", req.Username, "error", err)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генерация JWT с реальными данными пользователя
	accessToken, err := api.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		api.logger.Error("Failed to generate access token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := auth.LoginResponse{
		AccessToken: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

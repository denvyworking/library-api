package api

import (
	"encoding/json"
	"leti/pkg/auth"
	"net/http"
	"strings"
	"time"
)

// Login handles user authentication
// @Summary Аутентификация пользователя
// @Description Возвращает JWT токен и refresh токен при успешной аутентификации
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body auth.LoginRequest true "Учётные данные"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} string "Невалидные данные"
// @Failure 401 {object} string "Неверные учётные данные"
// @Router /api/v1/auth/login [post]
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
		return
	}

	user, err := api.srv.ValidateUserCredentials(r.Context(), req.Username, req.Password)
	if err != nil {
		api.logger.Error("Invalid credentials", "username", req.Username, "error", err)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := api.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		api.logger.Error("Failed to generate access token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := api.jwtService.GenerateRefreshToken()
	if err != nil {
		api.logger.Error("Failed to generate refresh token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(api.jwtService.GetRefreshTokenTTL())
	if err := api.srv.SaveRefreshToken(r.Context(), user.ID, refreshToken, expiresAt); err != nil {
		api.logger.Error("Failed to save refresh token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := auth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register handles user registration
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает токены
// @Tags auth
// @Accept json
// @Produce json
// @Param user body auth.RegisterRequest true "Данные пользователя"
// @Success 201 {object} auth.LoginResponse
// @Failure 400 {object} string "Невалидные данные"
// @Failure 409 {object} string "Пользователь уже существует"
// @Router /api/v1/auth/register [post]
func (api *api) register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
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
	if len(strings.TrimSpace(req.Username)) < 3 {
		http.Error(w, "username must be at least 3 characters", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		http.Error(w, "password cannot be empty", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	userID, err := api.srv.RegisterUser(r.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}
		api.logger.Error("Failed to register user", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	accessToken, err := api.jwtService.GenerateAccessToken(userID, req.Role)
	if err != nil {
		api.logger.Error("Failed to generate access token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := api.jwtService.GenerateRefreshToken()
	if err != nil {
		api.logger.Error("Failed to generate refresh token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(api.jwtService.GetRefreshTokenTTL())
	if err := api.srv.SaveRefreshToken(r.Context(), userID, refreshToken, expiresAt); err != nil {
		api.logger.Error("Failed to save refresh token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := auth.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Refresh handles token refresh
// @Summary Обновление токенов
// @Description Обновляет пару access/refresh токенов
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body auth.RefreshTokenRequest true "Refresh токен"
// @Success 200 {object} auth.RefreshTokenResponse
// @Failure 400 {object} string "Невалидные данные"
// @Failure 401 {object} string "Невалидный или истекший refresh токен"
// @Router /api/v1/auth/refresh [post]
func (api *api) refresh(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(req.RefreshToken) == "" {
		http.Error(w, "refresh_token cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := api.srv.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		api.logger.Error("Invalid refresh token", "error", err)
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	if err := api.srv.DeleteRefreshToken(r.Context(), req.RefreshToken); err != nil {
		api.logger.Warn("Failed to delete old refresh token", "error", err)
	}

	accessToken, err := api.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		api.logger.Error("Failed to generate access token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := api.jwtService.GenerateRefreshToken()
	if err != nil {
		api.logger.Error("Failed to generate refresh token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(api.jwtService.GetRefreshTokenTTL())
	if err := api.srv.SaveRefreshToken(r.Context(), user.ID, newRefreshToken, expiresAt); err != nil {
		api.logger.Error("Failed to save refresh token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := auth.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

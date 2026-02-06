package api

import (
	"encoding/json"
	"leti/pkg/api/dto"
	"leti/pkg/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_JWTAuthorization(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Создаём тестового пользователя (если нужно)
	// Но у нас уже есть Den из миграций

	t.Run("successful login", func(t *testing.T) {
		token := login(t, ts.URL, "Den", "password")
		require.NotEmpty(t, token)
	})

	t.Run("debug - check users", func(t *testing.T) {
		token := login(t, ts.URL, "Den", "password")
		require.NotEmpty(t, token)
	})
	t.Run("create book with valid token", func(t *testing.T) {
		token := login(t, ts.URL, "Den", "password")
		authorID := createAuthor(t, ts.URL, "Толстой_unique")
		genreID := createGenre(t, ts.URL, "Роман_unique")

		bookID := createBook(t, ts.URL, token, dto.CreateBookRequest{
			Name:     "Война и мир",
			AuthorID: authorID,
			GenreID:  genreID,
			Price:    1000,
		})
		require.Greater(t, bookID, 0)
	})

	t.Run("create book without token", func(t *testing.T) {
		authorID := createAuthor(t, ts.URL, "Достоевский")
		genreID := createGenre(t, ts.URL, "Роман")

		body := marshal(t, dto.CreateBookRequest{
			Name:     "Преступление и наказание",
			AuthorID: authorID,
			GenreID:  genreID,
			Price:    800,
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/books", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("create book with invalid token", func(t *testing.T) {
		authorID := createAuthor(t, ts.URL, "Чехов")
		genreID := createGenre(t, ts.URL, "Рассказ")

		body := marshal(t, dto.CreateBookRequest{
			Name:     "Вишнёвый сад",
			AuthorID: authorID,
			GenreID:  genreID,
			Price:    600,
		})

		req := newRequestWithAuth(t, http.MethodPost, ts.URL+"/api/books", "invalid-token", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("login returns refresh token", func(t *testing.T) {
		body := marshal(t, map[string]string{
			"username": "Den",
			"password": "password",
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/auth/login", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		require.NotEmpty(t, result["access_token"])
		require.NotEmpty(t, result["refresh_token"]) // должен быть refresh токен
	})
}

func TestE2E_Register(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("successful registration", func(t *testing.T) {
		body := marshal(t, map[string]string{
			"username": "newuser",
			"password": "password123",
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/register", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		require.NotEmpty(t, result["access_token"])
		require.NotEmpty(t, result["refresh_token"])
	})

	t.Run("registration with existing username", func(t *testing.T) {
		body := marshal(t, map[string]string{
			"username": "Den", // уже существует
			"password": "password123",
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/register", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("registration with invalid username (too short)", func(t *testing.T) {
		body := marshal(t, map[string]string{
			"username": "ab", // слишком короткий
			"password": "password123",
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/register", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("registration with invalid password (too short)", func(t *testing.T) {
		body := marshal(t, map[string]string{
			"username": "validuser",
			"password": "12345", // слишком короткий
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/register", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("registration with admin role", func(t *testing.T) {
		body := marshal(t, map[string]string{
			"username": "adminuser",
			"password": "password123",
			"role":     "admin",
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/register", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Проверяем, что можем логиниться с новым пользователем
		token := login(t, ts.URL, "adminuser", "password123")
		require.NotEmpty(t, token)
	})
}

func TestE2E_Refresh(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("successful refresh", func(t *testing.T) {
		// Сначала логинимся, чтобы получить refresh токен
		loginBody := marshal(t, map[string]string{
			"username": "Den",
			"password": "password",
		})

		loginReq := newRequest(t, http.MethodPost, ts.URL+"/api/auth/login", loginBody)
		loginResp := doRequest(t, loginReq)
		defer loginResp.Body.Close()

		var loginResult map[string]string
		json.NewDecoder(loginResp.Body).Decode(&loginResult)
		refreshToken := loginResult["refresh_token"]
		require.NotEmpty(t, refreshToken)

		// Теперь используем refresh токен
		refreshBody := marshal(t, map[string]string{
			"refresh_token": refreshToken,
		})

		refreshReq := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/refresh", refreshBody)
		refreshResp := doRequest(t, refreshReq)
		defer refreshResp.Body.Close()

		require.Equal(t, http.StatusOK, refreshResp.StatusCode)

		var refreshResult map[string]string
		json.NewDecoder(refreshResp.Body).Decode(&refreshResult)
		require.NotEmpty(t, refreshResult["access_token"])
		require.NotEmpty(t, refreshResult["refresh_token"])
		require.NotEqual(t, refreshToken, refreshResult["refresh_token"]) // новый токен должен отличаться
	})

	t.Run("refresh with invalid token", func(t *testing.T) {
		body := marshal(t, map[string]string{
			"refresh_token": "invalid-refresh-token",
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/refresh", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("refresh with empty token", func(t *testing.T) {
		body := marshal(t, map[string]string{
			"refresh_token": "",
		})

		req := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/refresh", body)
		resp := doRequest(t, req)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("refresh token rotation", func(t *testing.T) {
		// Логинимся
		loginBody := marshal(t, map[string]string{
			"username": "Den",
			"password": "password",
		})

		loginReq := newRequest(t, http.MethodPost, ts.URL+"/api/auth/login", loginBody)
		loginResp := doRequest(t, loginReq)
		defer loginResp.Body.Close()

		var loginResult map[string]string
		json.NewDecoder(loginResp.Body).Decode(&loginResult)
		oldRefreshToken := loginResult["refresh_token"]

		// Первый refresh
		refreshBody1 := marshal(t, map[string]string{
			"refresh_token": oldRefreshToken,
		})

		refreshReq1 := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/refresh", refreshBody1)
		refreshResp1 := doRequest(t, refreshReq1)
		defer refreshResp1.Body.Close()

		require.Equal(t, http.StatusOK, refreshResp1.StatusCode)

		var refreshResult1 map[string]string
		json.NewDecoder(refreshResp1.Body).Decode(&refreshResult1)
		newRefreshToken := refreshResult1["refresh_token"]

		// Старый токен больше не должен работать
		refreshBody2 := marshal(t, map[string]string{
			"refresh_token": oldRefreshToken,
		})

		refreshReq2 := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/refresh", refreshBody2)
		refreshResp2 := doRequest(t, refreshReq2)
		defer refreshResp2.Body.Close()

		require.Equal(t, http.StatusUnauthorized, refreshResp2.StatusCode) // старый токен инвалидирован

		// Новый токен должен работать
		refreshBody3 := marshal(t, map[string]string{
			"refresh_token": newRefreshToken,
		})

		refreshReq3 := newRequest(t, http.MethodPost, ts.URL+"/api/v1/auth/refresh", refreshBody3)
		refreshResp3 := doRequest(t, refreshReq3)
		defer refreshResp3.Body.Close()

		require.Equal(t, http.StatusOK, refreshResp3.StatusCode)
	})
}

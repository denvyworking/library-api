// @title Library API
// @version 1.0
// @description Production-ready book catalog API
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"leti/pkg/api"
	psg "leti/pkg/repository/postgres"
	"leti/pkg/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "leti/docs" // ← генерируется swag init

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func getDBConnectionString() string {
	if conn := os.Getenv("DB_CONNECTION_STRING"); conn != "" {
		return conn
	}
	return "postgresql://postgres:postgres@localhost:5432/courses?sslmode=disable"
}

func main() {
	connStr := getDBConnectionString()
	db, err := psg.New(connStr)
	if err != nil {
		slog.Error("Failed to connect to DB", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		authToken = "adminToken"
	}

	srv := service.NewService(db)
	router := mux.NewRouter()
	logger := slog.Default()
	apiHandler := api.New(router, srv, logger, authToken)
	apiHandler.RegistreRoutes()

	// Добавляем Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}

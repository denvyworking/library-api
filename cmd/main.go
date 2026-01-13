package main

import (
	"context"
	"leti/pkg/api"
	"leti/pkg/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	//"leti/pkg/models"
	psg "leti/pkg/repository/postgres"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
)

func main() {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@postgres:5432/courses"
	}

	db, err := psg.New(connStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		authToken = "adminToken"
	}

	srv := service.NewService(db)
	router := mux.NewRouter()
	logger := slog.Default()
	api := api.New(router, srv, logger, authToken)
	api.RegistreRoutes()

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
	// Docker и Kubernetes убьют процессы через 30 секунд по умолчанию
	// *** в данном проекте не целесообразно столько ждать, поэтому 10 сек. ***
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// getDBConnectionString возвращает строку подключения к БД.
func getDBConnectionString() string {
	if conn := os.Getenv("DB_CONNECTION_STRING"); conn != "" {
		return conn
	}
	// Для локальной разработки
	return "postgresql://postgres:postgres@localhost:5432/courses?sslmode=disable"
}

func getMigrationsPath() string {
	// Если задана переменная окружения - используем её (для Docker)
	if path := os.Getenv("MIGRATIONS_PATH"); path != "" {
		return "file://" + filepath.ToSlash(path)
	}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	// Поднимаемся: cmd/migrate → корень проекта → migrations/
	migrationsDir := filepath.Join(dir, "..", "..", "migrations")
	return "file://" + filepath.ToSlash(migrationsDir)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up|down|version]")
		os.Exit(1)
	}

	action := os.Args[1]
	connStr := getDBConnectionString()
	migrationsPath := getMigrationsPath()

	m, err := migrate.New(migrationsPath, connStr)
	if err != nil {
		fmt.Printf("Error creating migrator: %v\n", err)
		os.Exit(1)
	}

	switch action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("Error applying migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Printf("Error rolling back migrations: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("🔽 Migrations rolled back")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			fmt.Printf("Error getting version: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Version: %d, Dirty: %t\n", version, dirty)

	default:
		fmt.Printf("Unknown action: %s\n", action)
		os.Exit(1)
	}
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Set timezone to UTC for the entire program
	time.Local = time.UTC

	// Load .env.local
	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: .env.local not found: %v", err)
	}

	// Build DATABASE_URL from individual components
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal(
			"Required database environment variables are not set (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)",
		)
	}

	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=UTC",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	// Create migration instance
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Parse command from arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up completed successfully")

	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down completed successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Current version: %d (dirty: %v)\n", version, dirty)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migrate/main.go [command]")
	fmt.Println("\nCommands:")
	fmt.Println("  up       - Run all pending migrations")
	fmt.Println("  down     - Rollback last migration")
	fmt.Println("  version  - Show current migration version")
}

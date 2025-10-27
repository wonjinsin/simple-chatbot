package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Port       string
	Env        string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads configuration from .env.local file and environment variables
// Environment variables take priority over file values
func Load() *Config {
	// Try to load .env.local file (ignore error if file doesn't exist)
	_ = godotenv.Load(".env.local")

	cfg := &Config{
		Port:       mustGetEnv("PORT"),
		Env:        mustGetEnv("ENV"),
		DBHost:     mustGetEnv("DB_HOST"),
		DBPort:     mustGetEnv("DB_PORT"),
		DBUser:     mustGetEnv("DB_USER"),
		DBPassword: mustGetEnv("DB_PASSWORD"),
		DBName:     mustGetEnv("DB_NAME"),
		DBSSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}

	log.Printf("Configuration loaded: ENV=%s, PORT=%s, DB=%s@%s:%s/%s",
		cfg.Env, cfg.Port, cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	return cfg
}

// mustGetEnv reads an environment variable or panics if not found
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// getEnvOrDefault reads an environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDatabaseURL constructs PostgreSQL connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=UTC",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

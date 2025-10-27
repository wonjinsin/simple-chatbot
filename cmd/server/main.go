package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wonjinsin/go-boilerplate/internal/config"
	httpHandler "github.com/wonjinsin/go-boilerplate/internal/handler/http"
	"github.com/wonjinsin/go-boilerplate/internal/repository/postgres"
	"github.com/wonjinsin/go-boilerplate/internal/usecase"
	"github.com/wonjinsin/go-boilerplate/pkg/logger"
)

func main() {
	// Print ASCII art banner
	printBanner()

	// Set timezone to UTC for the entire program
	time.Local = time.UTC

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Initialize(cfg.Env)

	// Initialize database connection
	userRepo, err := postgres.NewUserRepository(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer func() {
		if closer, ok := userRepo.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				log.Printf("failed to close database: %v", err)
			}
		}
	}()

	// Wiring (Composition Root)
	var userSvc usecase.UserService = usecase.NewUserService(userRepo)

	// Create chi router
	router := httpHandler.NewRouter(userSvc)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("HTTP server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		_ = srv.Close()
	}
	log.Println("bye")
}

func printBanner() {
	// Read banner from file
	bannerPath := "internal/config/banner.asc"
	bannerBytes, err := os.ReadFile(bannerPath)
	if err != nil {
		log.Printf("warning: could not read banner file: %v", err)
		return
	}

	log.Println(string(bannerBytes))
}

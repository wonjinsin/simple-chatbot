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

	"github.com/wonjinsin/simple-chatbot/internal/config"
	"github.com/wonjinsin/simple-chatbot/internal/database"
	httpHandler "github.com/wonjinsin/simple-chatbot/internal/handler/http"
	chatgptRepo "github.com/wonjinsin/simple-chatbot/internal/repository/langchain/chatGPT"
	langchain "github.com/wonjinsin/simple-chatbot/internal/repository/langchain/ollama"
	"github.com/wonjinsin/simple-chatbot/internal/repository/postgres"
	"github.com/wonjinsin/simple-chatbot/internal/usecase"
	"github.com/wonjinsin/simple-chatbot/pkg/logger"
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

	// Initialize LLM
	ollamaLLM, err := database.NewOllamaLLM()
	if err != nil {
		log.Fatalf("failed to initialize LLM: %v", err)
	}

	// Initialize ChatGPT Embedder
	chatGPTEmbedder, err := database.NewChatGPTEmbedder(cfg.OpenAIAPIKey)
	if err != nil {
		log.Fatalf("failed to initialize ChatGPT embedder: %v", err)
	}

	// Initialize PostgreSQL database connection
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	basicChatRepo := langchain.NewBasicChatRepo(ollamaLLM)
	embeddingRepo := chatgptRepo.NewEmbeddingRepository(chatGPTEmbedder)
	inquiryKnowledgeRepo := postgres.NewInquiryKnowledgeRepository(db)

	// Wiring (Composition Root)
	userSvc := usecase.NewUserService(userRepo)
	basicChatSvc := usecase.NewBasicChatService(basicChatRepo)
	inquirySvc := usecase.NewInquiryServiceImpl(embeddingRepo, inquiryKnowledgeRepo)

	// Create chi router
	router := httpHandler.NewRouter(userSvc, basicChatSvc, inquirySvc)
	handler := http.TimeoutHandler(router, 59*time.Second, "Timeout")

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
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

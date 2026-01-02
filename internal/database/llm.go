package database

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/model/ollama"
)

func NewOllamaLLM() (*ollama.ChatModel, error) {
	ctx := context.Background()
	return ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		// Basic Configuration
		BaseURL: "http://localhost:11434", // Ollama service address
		Timeout: 30 * time.Second,         // Request timeout

		// Model Configuration
		Model: "gemma3:1b", // Model name
	})
}

func NewChatGPTEmbedder(k string) (*openai.Embedder, error) {
	ctx := context.Background()
	return openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		APIKey:  k,
		Model:   "text-embedding-3-small",
		Timeout: 30 * time.Second,
	})
}

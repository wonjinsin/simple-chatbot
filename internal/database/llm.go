package database

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

func NewOllamaLLM() (*ollama.ChatModel, error) {
	ctx := context.Background()
	model, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		// Basic Configuration
		BaseURL: "http://localhost:11434", // Ollama service address
		Timeout: 30 * time.Second,         // Request timeout

		// Model Configuration
		Model: "gemma3:1b", // Model name
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ollama chat model")
	}
	return model, nil
}

func NewChatGPTEmbedder(k string) (*openai.Embedder, error) {
	ctx := context.Background()
	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		APIKey:  k,
		Model:   "text-embedding-3-small",
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create openai embedder")
	}
	return embedder, nil
}

package langchain

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

type embeddingRepo struct {
	embedder *openai.Embedder
}

// NewEmbeddingRepository creates a new embedding repository
func NewEmbeddingRepository(embedder *openai.Embedder) repository.EmbeddingRepository {
	return &embeddingRepo{embedder: embedder}
}

// EmbedStrings converts text strings to embedding vectors using ChatGPT
func (r *embeddingRepo) EmbedStrings(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings, err := r.embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to embed strings")
	}
	return embeddings, nil
}

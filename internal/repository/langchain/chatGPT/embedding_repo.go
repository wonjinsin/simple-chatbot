package langchain

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/internal/domain"
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

// EmbedString converts text string to embedding vector using ChatGPT
func (r *embeddingRepo) EmbedString(
	ctx context.Context,
	text string,
) (domain.Embedding, error) {
	embeddings, err := r.embedder.EmbedStrings(ctx, []string{text})
	if err != nil {
		return nil, errors.Wrap(err, "failed to embed string")
	}
	if len(embeddings) == 0 {
		return nil, errors.New(
			constants.InternalError,
			"embedding generation returned empty result",
			nil,
		)
	}
	return domain.NewEmbedding(embeddings[0]), nil
}

// EmbedStrings converts text strings to embedding vectors using ChatGPT
func (r *embeddingRepo) EmbedStrings(
	ctx context.Context,
	texts []string,
) (domain.Embeddings, error) {
	embeddings, err := r.embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to embed strings")
	}
	return domain.NewEmbeddings(embeddings), nil
}

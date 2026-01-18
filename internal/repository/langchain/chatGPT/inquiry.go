package langchain

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

type inquiryRepo struct {
	embedder *openai.Embedder
}

// NewInquiryRepo creates a new inquiry repository
func NewInquiryRepo(embedder *openai.Embedder) repository.InquiryRepository {
	return &inquiryRepo{embedder: embedder}
}

// EmbedStrings embeds strings and returns embeddings
func (r *inquiryRepo) EmbedStrings(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings, err := r.embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to embed strings")
	}
	return embeddings, nil
}

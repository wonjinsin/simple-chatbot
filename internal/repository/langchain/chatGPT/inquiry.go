package langchain

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
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
	return r.embedder.EmbedStrings(ctx, texts)
}

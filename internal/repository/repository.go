package repository

import (
	"context"

	"github.com/wonjinsin/simple-chatbot/internal/domain"
)

// AnswerRefineRepository defines the interface for refining answers based on context
type AnswerRefineRepository interface {
	RefineAnswer(ctx context.Context, contextStr string) (string, error)
}

// EmbeddingRepository defines the interface for text embedding operations
type EmbeddingRepository interface {
	// EmbedString converts text string to embedding vector using LLM
	EmbedString(ctx context.Context, text string) (domain.Embedding, error)
	// EmbedStrings converts text strings to embedding vectors using LLM
	EmbedStrings(ctx context.Context, texts []string) (domain.Embeddings, error)
}

// InquiryKnowledgeRepository defines the interface for inquiry knowledge database operations
type InquiryKnowledgeRepository interface {
	// BatchSaveInquiryKnowledge saves multiple inquiry knowledge entries to database
	BatchSaveInquiryKnowledge(ctx context.Context, items domain.InquiryKnowledges) error
	// FindSimilar finds inquiry knowledge entries similar to the given embedding vector with
	// similarity scores
	FindSimilars(
		ctx context.Context,
		embedding domain.Embedding,
		limit int,
	) (domain.InquirySimilarityResults, error)
}

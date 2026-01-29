package repository

import (
	"context"

	"github.com/wonjinsin/simple-chatbot/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Save(*domain.User) error
	FindByID(id int) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	List(offset, limit int) (domain.Users, error)
}

// BasicChatRepository defines the interface for basic chat data access
type BasicChatRepository interface {
	AskBasicChat(ctx context.Context, msg string) (string, error)
	AskBasicPromptTemplateChat(ctx context.Context, msg string) (string, error)
	AskBasicParallelChat(ctx context.Context, msg string) (string, error)
	AskBasicBranchChat(ctx context.Context, msg string) (string, error)
	AskWithTool(ctx context.Context, msg string) (string, error)
	AskWithGraph(ctx context.Context, msg string) (string, error)
	AskWithGraphWithBranch(ctx context.Context, _ string) (string, error)
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

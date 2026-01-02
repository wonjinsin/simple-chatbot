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

// InquiryRepository defines the interface for inquiry data access
type InquiryRepository interface {
	EmbedStrings(ctx context.Context, texts []string) ([][]float64, error)
}

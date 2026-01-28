package usecase

import (
	"context"

	"github.com/wonjinsin/simple-chatbot/internal/domain"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, name, email string) (*domain.User, error)
	GetUser(ctx context.Context, id int) (*domain.User, error)
	ListUsers(ctx context.Context, offset, limit int) (domain.Users, error)
}

// BasicChatService defines the interface for basic chat business logic
type BasicChatService interface {
	AskBasicChat(ctx context.Context, msg string) (string, error)
	AskBasicPromptTemplateChat(ctx context.Context, msg string) (string, error)
}

// InquiryService defines the interface for inquiry business logic
type InquiryService interface {
	Ask(ctx context.Context, msg string) (*domain.InquirySimilarityResult, error)
	EmbedInquiryOrigins(ctx context.Context) error
}

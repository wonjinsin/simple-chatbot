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
	Ask(ctx context.Context, msg string) (string, error)
}

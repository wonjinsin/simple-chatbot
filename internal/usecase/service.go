package usecase

import (
	"context"

	"github.com/wonjinsin/go-boilerplate/internal/domain"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, name, email string) (*domain.User, error)
	GetUser(ctx context.Context, id int) (*domain.User, error)
	ListUsers(ctx context.Context, offset, limit int) (domain.Users, error)
}

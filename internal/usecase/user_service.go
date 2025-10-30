package usecase

import (
	"context"
	"time"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/internal/domain"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	// Check if user with email already exists
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		// If error is NotFound, it's okay - user doesn't exist yet
		if !errors.HasCode(err, constants.NotFound) {
			return nil, errors.Wrap(err, "failed to check existing email")
		}
	} else if existing != nil {
		// User exists - duplicate email
		return nil, errors.New(constants.ConstraintError, "duplicate email", nil)
	}

	// ID is 0 - database will auto-generate
	u, err := domain.NewUser(0, name, email, time.Now())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}
	if err := s.repo.Save(u); err != nil {
		return nil, errors.Wrap(err, "failed to save user")
	}
	return u, nil
}

func (s *userService) GetUser(ctx context.Context, id int) (*domain.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	return u, nil
}

func (s *userService) ListUsers(ctx context.Context, offset, limit int) (domain.Users, error) {
	if limit <= 0 {
		limit = 50
	}
	users, err := s.repo.List(offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}
	return users, nil
}

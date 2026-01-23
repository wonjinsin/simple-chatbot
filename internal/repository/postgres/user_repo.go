package postgres

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/internal/domain"
	"github.com/wonjinsin/simple-chatbot/internal/repository"
	"github.com/wonjinsin/simple-chatbot/internal/repository/postgres/dao/ent"
	"github.com/wonjinsin/simple-chatbot/internal/repository/postgres/dao/ent/user"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

type userRepo struct {
	client *ent.Client
}

// NewUserRepository creates a new PostgreSQL-based user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	// Create ent client
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	return &userRepo{client: client}
}

// Close closes the database connection
func (r *userRepo) Close() error {
	if err := r.client.Close(); err != nil {
		return errors.Wrap(err, "failed to close database connection")
	}
	return nil
}

// Save creates or updates a user
func (r *userRepo) Save(u *domain.User) error {
	ctx := context.Background()

	// Apply transformations using mapper
	name, email := toEntUserData(u)

	// Check if user already exists
	if u.ID != 0 {
		// Update existing user
		_, err := r.client.User.
			UpdateOneID(u.ID).
			SetName(name).
			SetEmail(email).
			Save(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to update user")
		}
		return nil
	}

	// Create new user
	created, err := r.client.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetCreatedAt(u.CreatedAt).
		Save(ctx)
	if err != nil {
		// Check for duplicate email error
		if ent.IsConstraintError(err) {
			return errors.New(constants.ConstraintError, "duplicate email", err)
		}
		return errors.Wrap(err, "failed to create user")
	}

	// Update domain object with generated ID
	u.ID = created.ID
	return nil
}

// FindByID retrieves a user by ID
func (r *userRepo) FindByID(id int) (*domain.User, error) {
	ctx := context.Background()

	u, err := r.client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New(constants.NotFound, "user not found", err)
		}
		return nil, errors.Wrap(err, "failed to find user")
	}

	return toDomainUser(u), nil
}

// FindByEmail retrieves a user by email
func (r *userRepo) FindByEmail(email string) (*domain.User, error) {
	ctx := context.Background()

	u, err := r.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New(constants.NotFound, "user not found", err)
		}
		return nil, errors.Wrap(err, "failed to find user by email")
	}

	return toDomainUser(u), nil
}

// List retrieves a list of users with pagination
func (r *userRepo) List(offset, limit int) (domain.Users, error) {
	ctx := context.Background()

	users, err := r.client.User.
		Query().
		Order(ent.Asc(user.FieldCreatedAt)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}

	result := make(domain.Users, len(users))
	for i, u := range users {
		result[i] = toDomainUser(u)
	}

	return result, nil
}

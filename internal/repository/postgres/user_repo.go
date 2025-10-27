package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/wonjinsin/go-boilerplate/internal/config"
	"github.com/wonjinsin/go-boilerplate/internal/constants"
	"github.com/wonjinsin/go-boilerplate/internal/domain"
	"github.com/wonjinsin/go-boilerplate/internal/repository"
	"github.com/wonjinsin/go-boilerplate/internal/repository/postgres/dao/ent"
	"github.com/wonjinsin/go-boilerplate/internal/repository/postgres/dao/ent/user"
	"github.com/wonjinsin/go-boilerplate/pkg/errors"
)

type userRepo struct {
	client *ent.Client
}

// NewUserRepository creates a new PostgreSQL-based user repository
func NewUserRepository(cfg *config.Config) (repository.UserRepository, error) {
	// Open database connection
	db, err := sql.Open("pgx", cfg.GetDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create ent client
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	return &userRepo{client: client}, nil
}

// Close closes the database connection
func (r *userRepo) Close() error {
	return r.client.Close()
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
			return fmt.Errorf("failed to update user: %w", err)
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

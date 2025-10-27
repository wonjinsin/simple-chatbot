package postgres

import (
	"github.com/wonjinsin/go-boilerplate/internal/domain"
	"github.com/wonjinsin/go-boilerplate/internal/repository/postgres/dao/ent"
)

// toDomainUser converts ent.User to domain.User
func toDomainUser(u *ent.User) *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// toEntUserData converts domain.User to ent values for creation/update
// This is where you can add transformation logic if needed
func toEntUserData(u *domain.User) (name string, email string) {
	return u.Name, u.Email
}

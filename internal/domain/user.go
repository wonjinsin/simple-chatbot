package domain

import (
	"time"

	"github.com/wonjinsin/go-boilerplate/internal/constants"
	pkgConstants "github.com/wonjinsin/go-boilerplate/pkg/constants"
	"github.com/wonjinsin/go-boilerplate/pkg/errors"
	"github.com/wonjinsin/go-boilerplate/pkg/utils"
)

// User is an aggregate root.
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(id int, name, email string, now time.Time) (*User, error) {
	name = utils.NormalizeName(name)
	email = utils.NormalizeEmail(email)

	if utils.IsEmptyOrWhitespace(name) || len(name) > pkgConstants.MaxNameLength {
		return nil, errors.New(constants.InvalidParameter, "invalid name", nil)
	}
	if !utils.IsValidEmail(email) {
		return nil, errors.New(constants.InvalidParameter, "invalid email format", nil)
	}
	return &User{ID: id, Name: name, Email: email, CreatedAt: now}, nil
}

type Users []*User

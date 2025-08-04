package domain

import (
	"context"
	"time"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	UserID         string
	Name           *string // for oauth2 user
	Username       *string
	FirstName      *string
	LastName       *string
	Bio            *string
	ProfilePicture *string
	IsVerified     bool
	Email          string
	Password       *string
	RefreshToken   *string
	AccessToken    *string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Provider string // for oauth2 user
	RawData  map[string]any

	Role Role
}

// UserRepository Interface
type IUserRepository interface {
	CreateUser(c context.Context, user *User) error
	FindByEmail(c context.Context, email string) (*User, error)
	FindByID(c context.Context, id string) (*User, error)
	UpdateUser(c context.Context, user *User) error
	DeleteByID(c context.Context, id string) error

	SetEmailVerified(c context.Context, userID string) error
	IsEmailVerified(c context.Context, userID string) (bool, error)

	CountByEmail(c context.Context, email string) (int64, error)
	CountAll(c context.Context) (int64, error)
}

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

//PasswordService Interface

type IPasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, inputPassword string) bool
}

//JWTService Interface

type IJWTService interface {
	GenerateVerificationToken(userID string) (string, error)
	GenerateAccessToken(userID string, role string) (string, error)
	GenerateRefreshToken(userID string, role string) (string, error)
	ValidateToken(token string, isRefresh bool) (string, error)
}

type OAuth2ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// OAuth2 providers interface
type IOAuth2Provider interface {
	Name() string // provider name
	Authenticate(ctx context.Context, code string) (*User, error)
	GetAuthorizationURL(state string) string
}

type IOAuth2Service interface {
	SupportedProviders() []string
	GetAuthorizationURL(provider string, state string) (string, error)
	Authenticate(ctx context.Context, provider string, code string) (*User, error)
}

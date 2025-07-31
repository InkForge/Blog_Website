package domain

import "time"

type User struct {
	UserID         int
	RoleID         int
	OAuthID        *int
	Username       *string
	FirstName      *string
	LastName       *string
	Bio            *string
	ProfilePicture *string
	Email          string
	HashedPassword *string
	RefreshToken   *string
	AccessToken    *string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Role *Role
}

type Role struct {
	RoleID int
	Role   string
}

type OAuthUser struct {
	ID           int
	ProviderName string
	ProviderID   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

//UserRepository Interface

type IUserRepository interface {
	CreateUser(user User) error
	FindByEmail(email string) (User, error)
	CountByEmail(email string) (int64, error)
	CountAll() (int64, error)
}

//PasswordService Interface

type IPasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, inputPassword string) bool
}

//JWTService Interface

type IJWTService interface {
	GenerateAccessToken(userID string, role string) (string, error)
	GenerateRefreshToken(userID string, role string) (string, error)
	ValidateToken(token string, isRefresh bool) (string, error)
}

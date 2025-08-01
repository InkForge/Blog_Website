package domain

import "time"

type Role string

const (
  RoleUser     Role = "USER"
  RoleAdmin    Role = "ADMIN"
  
)


type User struct {
	UserID         string
	
	Username       *string
	FirstName      *string
	LastName       *string
	Bio            *string
	ProfilePicture *string
	IsVerified 		bool
	Email          string
	Password     *string
	RefreshToken   *string
	AccessToken    *string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Role     Role
}


//UserRepository Interface

type IUserRepository interface {
	CreateUser(user User) error
	FindByEmail(email string) (User, error)
	CountByEmail(email string) (int64, error)
	CountAll() (int64, error)
	IsVerified(email string)bool
}

//PasswordService Interface

type IPasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, inputPassword string) bool
}

//JWTService Interface

type IJWTService interface {
	GenerateVerificationToken(userID string)(string,error)
	GenerateAccessToken(userID string, role string) (string, error)
	GenerateRefreshToken(userID string, role string) (string, error)
	ValidateRefreshToken(token string) (userID string, role string, err error)
	ValidateAccessToken(token string) (userID string, role string, err error)
	ValidateVerificationToken(token string) (userID string, err error)
	GeneratePasswordResetToken(userID string) (string, error)
	ValidatePasswordResetToken(token string) (userID string, err error)
	RevokeRefreshToken(token string) error
	IsRefreshTokenRevoked(token string) (bool, error)

		
}

type IRevocationRepository interface {
    RevokeRefreshToken(token string, expiresAt time.Time) error
    IsRefreshTokenRevoked(token string) (bool, error)
}
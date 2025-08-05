package domain

import (
	"context"
	"time"
)

type IAuthUsecase interface {
	// Register registers a new user using provided user data.
	// It should validate the input, hash the password, store the user in the database,
	// and send a verification email.
	Register(ctx context.Context, user, oauthUser *User) (*User, error)
	
	// Login authenticates a user using an identifier (username or email) and password.
	// It should verify credentials, check if the email is verified, and return accesstoken, refreshtoken, user data.
	Login(ctx context.Context, input *User) (string, string, *User, error)

	// Logout logs out a user by invalidating their session or deleting the stored refresh token.
	// This ensures the user can no longer refresh their access token.
	Logout(ctx context.Context, userID string) error

	// RefreshToken validates the provided refresh token and issues a new access token.
	// It should also return a new refresh token (optional depending on rotation strategy),
	// and indicate how long the access token is valid for.
	// The refresh token is stored in the user's DB record.
	// returns access, refresh, duration the until access token dies, error
	RefreshToken(ctx context.Context, userID string) (*string, *string, time.Duration, error)

	// VerifyEmail verifies the user's email address using a token sent via email.
	// It should mark the email as verified in the database if the token is valid.
	VerifyEmail(ctx context.Context, token string) error

	// ResendVerificationEmail re-sends the email verification token to the user's email.
	// Should be used if the user didn’t receive or lost the initial verification email.
	ResendVerificationEmail(ctx context.Context, email string) error

	// RequestPasswordReset initiates a password reset flow by sending a reset token
	// to the user's email. The token is typically time-limited and signed.
	RequestPasswordReset(ctx context.Context, email string) error

	// ResetPassword resets the user's password using the provided reset token.
	// The token should be validated, and the new password hashed and saved.
	ResetPassword(ctx context.Context, token, newPassword string) error

	// ChangePassword allows an authenticated user to change their password by
	// verifying the old password and updating with the new one.
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
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
	ValidateRefreshToken(token string) (userID string, role string, err error)
	ValidateAccessToken(token string) (userID string, role string, err error)
	ValidateVerificationToken(token string) (userID string, err error)
	GeneratePasswordResetToken(userID string) (string, error)
	ValidatePasswordResetToken(token string) (userID string, err error)
	RevokeRefreshToken(token string) error
	IsRefreshTokenRevoked(token string) (bool, error)
	GetAccessTokenRemaining(token string) (time.Duration, error)
}

type IRevocationRepository interface {
	RevokeRefreshToken(token string, expiresAt time.Time) error
	IsRefreshTokenRevoked(token string) (bool, error)
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

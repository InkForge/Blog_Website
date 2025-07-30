package infrastructures

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type IJWTService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string, isRefresh bool) (string, error)
}

type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewJWTService(accessSecret string, refreshSecret string) IJWTService {
	return &JWTService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (j *JWTService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Shorter expiry for access token
		"iat": time.Now().Unix(),
		"iss": "your-app-name",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.accessSecret)
}

func (j *JWTService) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // Longer expiry for refresh token
		"iat": time.Now().Unix(),
		"iss": "your-app-name",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshSecret)
}

func (j *JWTService) ValidateToken(tokenString string, isRefresh bool) (string, error) {
	secret := j.accessSecret
	if isRefresh {
		secret = j.refreshSecret
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if subject, ok := claims["sub"].(string); ok {
			return subject, nil
		}
		return "", errors.New("invalid subject in token")
	}

	return "", errors.New("invalid token")
}
package infrastructures

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	userRepo domain.IUserRepository
}

func NewJWTService(accessSecret string, refreshSecret string, userRepo domain.IUserRepository)  domain.IJWTService {
	return &JWTService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		userRepo: userRepo,
	}
}
//generate access token 
func (j *JWTService) GenerateAccessToken(userID string, role string) (string, time.Duration, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(15 * time.Minute).Unix(), // Shorter expiry for access token
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err :=  token.SignedString(j.accessSecret)
	if err != nil {
		return "", 0, err
	}
	return tokenString, 15 * time.Minute, nil
}

//generate refresh token with longer expiry time 

func (j *JWTService) GenerateRefreshToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,

		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), 
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshSecret)
}

//generate verification token 
func (j *JWTService) GenerateVerificationToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,

		"exp": time.Now().Add(15 * time.Minute).Unix(), // Shorter expiry for verification token
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.accessSecret)
}

// helper to parse and validate, returning claims map
func (j *JWTService) parseToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

//validate access token 

func (j *JWTService) ValidateAccessToken(tokenString string) (string, string, error) {
	claims, err := j.parseToken(tokenString, j.accessSecret)
	if err != nil {
		return "", "", err
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid subject in token")
	}
	role, _ := claims["role"].(string) 
	return sub, role, nil
}

func (j *JWTService) ValidateRefreshToken(tokenString string) (string, string, error) {
	claims, err := j.parseToken(tokenString, j.refreshSecret)
	if err != nil {
		return "", "", err
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid subject in token")
	}

	role, _ := claims["role"].(string)

	return sub, role, nil
}




func (j *JWTService) ValidateVerificationToken(tokenString string) (string, error) {
	claims, err := j.parseToken(tokenString, j.accessSecret)
	if err != nil {
		return "", err
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid subject in token")
	}
	return sub, nil
}

func (j *JWTService) GeneratePasswordResetToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":     userID,
		"exp":     time.Now().Add(30 * time.Minute).Unix(), 
		"iat":     time.Now().Unix(),
		"purpose": "password_reset",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.accessSecret)
}

func (j *JWTService) ValidatePasswordResetToken(tokenString string) (string, error) {
	claims, err := j.parseToken(tokenString, j.accessSecret)
	if err != nil {
		return "", err
	}
	if purpose, ok := claims["purpose"].(string); !ok || purpose != "password_reset" {
		return "", errors.New("invalid token purpose")
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid subject in token")
	}
	return sub, nil
}



// helper to extract exp claim as int64
func extractExp(claims jwt.MapClaims) (int64, error) {
	expVal, ok := claims["exp"]
	if !ok {
		return 0, errors.New("exp claim not present")
	}

	switch v := expVal.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, errors.New("invalid exp claim format")
		}
		return n, nil
	default:
		return 0, errors.New("unexpected exp claim type")
	}
}

// Generic: remaining duration for any token given its secret
func (j *JWTService) GetTokenRemaining(tokenString string, secret []byte) (time.Duration, error) {
	claims, err := j.parseToken(tokenString, secret)
	if err != nil {
		return 0, err
	}

	expUnix, err := extractExp(claims)
	if err != nil {
		return 0, err
	}

	expTime := time.Unix(expUnix, 0)
	remaining := time.Until(expTime)
	if remaining <= 0 {
		return 0, errors.New("token already expired")
	}
	return remaining, nil
}

// Specific: access token remaining duration
func (j *JWTService) GetAccessTokenRemaining(tokenString string) (time.Duration, error) {
	return j.GetTokenRemaining(tokenString, j.accessSecret)
}

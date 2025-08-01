package infrastructures

import (
	"errors"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	revocationRepo domain.IRevocationRepository
}

func NewJWTService(accessSecret string, refreshSecret string, revocationRepo domain.IRevocationRepository)  domain.IJWTService {
	return &JWTService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		revocationRepo: revocationRepo,
	}
}
//generate access token 
func (j *JWTService) GenerateAccessToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(15 * time.Minute).Unix(), // Shorter expiry for access token
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.accessSecret)
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

	//check if it is revoked or not 
	
	revoked, err := j.IsRefreshTokenRevoked(tokenString)
	if err != nil {
		return "", "", err
	}
	if revoked {
		return "", "", errors.New("refresh token revoked")
	}

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

func (j *JWTService) RevokeRefreshToken(tokenString string) error {
    claims, err := j.parseToken(tokenString, j.refreshSecret)
    if err != nil {
        return err
    }
    expFloat, ok := claims["exp"].(float64)
    if !ok {
        return errors.New("invalid exp claim")
    }
    exp := time.Unix(int64(expFloat), 0)
    return j.revocationRepo.RevokeRefreshToken(tokenString, exp)
}

//function to check if th refresh token is revocked 

func (j *JWTService) IsRefreshTokenRevoked(tokenString string) (bool, error) {
    return j.revocationRepo.IsRefreshTokenRevoked(tokenString)
}
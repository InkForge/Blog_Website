package models

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID         primitive.ObjectID `bson:"_id,omitempty"`
	Name           *string            `bson:"name"`
	Username       *string            `bson:"username"`
	FirstName      *string            `bson:"first_name"`
	LastName       *string            `bson:"last_name"`
	Bio            *string            `bson:"bio"`
	ProfilePicture *string            `bson:"profile_picture"`
	IsVerified     bool               `bson:"is_verified"`
	Email          string             `bson:"email"`
	Password       *string            `bson:"password"`
	RefreshToken   *string            `bson:"refresh_token"`
	AccessToken    *string            `bson:"access_token"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`

	Provider string         `bson:"provider"`
	RawData  map[string]any `bson:"raw_data"`

	Role string `bson:"role"`
}

func UserFromDomain(u domain.User) (*User, error) {
	var objID primitive.ObjectID
	var err error

	if u.UserID != "" {
		objID, err = primitive.ObjectIDFromHex(u.UserID)
		if err != nil {
			return nil, domain.ErrInvalidUserID
		}
	}

	return &User{
		UserID: objID,
		Name:   u.Name,

		Username:       u.Username,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Bio:            u.Bio,
		ProfilePicture: u.ProfilePicture,
		IsVerified:     u.IsVerified,
		Email:          u.Email,
		Password:       u.Password,
		RefreshToken:   u.RefreshToken,
		AccessToken:    u.AccessToken,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,

		Provider: u.Provider,
		RawData:  u.RawData,

		Role: string(u.Role),
	}, nil
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		UserID:         u.UserID.Hex(),
		Name:           u.Name,
		Username:       u.Username,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Bio:            u.Bio,
		ProfilePicture: u.ProfilePicture,
		IsVerified:     u.IsVerified,
		Email:          u.Email,
		Password:       u.Password,
		RefreshToken:   u.RefreshToken,
		AccessToken:    u.AccessToken,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,

		Provider: u.Provider,
		RawData:  u.RawData,

		Role: domain.Role(u.Role),
	}
}

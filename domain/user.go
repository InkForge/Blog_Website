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
  
  Role      *Role
  
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

package domain

import (
	"context"
	"time"
)

type BlogView struct {
	ID       string
	Blog_id  string
	User_id  string
	ViewedAt time.Time
}

type IBlogViewRepository interface {
	CreateViewRecord(ctx context.Context, blog_id, user_id string) error
}

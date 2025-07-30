package domain

import (
	"context"
	"time"
)

type Blog struct {
	Blog_id        string    `json:"blog_id"`
	User_id        string    `json:"user_id"`
	Title          string    `json:"title"`
	Images         []string  `json:"images"`
	Content        string    `json:"content"`
	Tag_ids        []string  `json:"tag_ids"`
	Comment_ids    []string  `json:"comment_ids"`
	Posted_at      time.Time `json:"posted_at"`
	Like_counts    int       `json:"like_counts"`
	Dislike_counts int       `json:"dislike_counts"`
	Share_count    int       `json:"share_count"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
}

type IBlogRepository interface {
	Create(ctx context.Context, blog Blog) (string, error)
	GetAll(ctx context.Context) ([]Blog, error)
	GetByID(ctx context.Context, blogID string) (Blog, error)
	Update(ctx context.Context, blog Blog) error
	Delete(ctx context.Context, blogID string) error
}


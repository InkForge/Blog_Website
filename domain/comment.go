package domain

import (
	"context"
	"time"
)

type Comment struct {
	Comment_id  string
	Blog_id     string
	User_id     string
	Content     string
	Created_at  time.Time
	Updated_at  time.Time
}

type ICommentRepository interface {
	Create(ctx context.Context, comment Comment) (string, error)
	GetByID(ctx context.Context, commentID string) (Comment, error)
	GetByBlogID(ctx context.Context, blogID string) ([]Comment, error)
	Update(ctx context.Context, comment Comment) error
	Delete(ctx context.Context, commentID string) error
} 
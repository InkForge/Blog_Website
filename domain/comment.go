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

type ICommentUsecase interface {
	AddComment(ctx context.Context, blogID string, comment *Comment) (string, error)
	RemoveComment(ctx context.Context, blogID, commentID string) error
	GetBlogComments(ctx context.Context, blogID string) ([]Comment, error)
	UpdateComment(ctx context.Context, commentID string, comment *Comment) error
} 
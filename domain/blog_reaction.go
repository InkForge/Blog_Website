package domain

import (
	"context"
	"time"
)

type BlogReaction struct {
	ID            string
	Blog_id       string
	User_id       string
	Reaction_type int
	Created_at    time.Time
}

type IBlogReactionRepository interface {
	//CRUD
	CreateReaction(ctx context.Context, blogReaction BlogReaction) error
	GetReactionByUserAndBlog(ctx context.Context, blog_id, user_id string) (BlogReaction, error)
	UpdateReaction(ctx context.Context, blogReaction BlogReaction) error

	DeleteReaction(ctx context.Context, blog_id, user_id string) error
}

type IBlogReactionUsecase interface {
	LikeBlog(ctx context.Context, blogID, userID string) error

	DislikeBlog(ctx context.Context, blogID, userID string) error

	UnlikeBlog(ctx context.Context, blogID, userID string) error

	UndislikeBlog(ctx context.Context, blogID, userID string) error
}

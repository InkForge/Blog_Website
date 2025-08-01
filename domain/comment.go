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
	Like        int
	Dislike     int
	Created_at  time.Time
	Updated_at  time.Time
}

type CommentReaction struct {
	Comment_id  string
	User_id     string
	Action      int // -1 for DISLIKE, 0 for NONE, 1 for LIKE
	Created_at  time.Time
	Updated_at  time.Time
}

type ICommentRepository interface {
	Create(ctx context.Context, comment Comment) (string, error)
	GetByID(ctx context.Context, commentID string) (Comment, error)
	GetByBlogID(ctx context.Context, blogID string) ([]Comment, error)
	Update(ctx context.Context, comment Comment) error
	Delete(ctx context.Context, commentID string) error
	UpdateReactionCounts(ctx context.Context, commentID string, likeCount, dislikeCount int) error
}

type ICommentReactionRepository interface {
	Create(ctx context.Context, reaction CommentReaction) error
	GetByCommentAndUser(ctx context.Context, commentID, userID string) (CommentReaction, error)
	Update(ctx context.Context, reaction CommentReaction) error
	Delete(ctx context.Context, commentID, userID string) error
	GetReactionCounts(ctx context.Context, commentID string) (int, int, error)
}

type ICommentUsecase interface {
	AddComment(ctx context.Context, blogID string, comment *Comment) (string, error)
	RemoveComment(ctx context.Context, blogID, commentID string) error
	GetBlogComments(ctx context.Context, blogID string) ([]Comment, error)
	UpdateComment(ctx context.Context, commentID string, comment *Comment) error
}

type ICommentReactionUsecase interface {
	LikeComment(ctx context.Context, commentID, userID string) error
	DislikeComment(ctx context.Context, commentID, userID string) error
	RemoveReaction(ctx context.Context, commentID, userID string) error
	GetUserReaction(ctx context.Context, commentID, userID string) (int, error)
} 
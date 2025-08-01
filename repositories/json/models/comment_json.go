package models

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type CommentJSON struct {
	CommentID string    `json:"comment_id"`
	BlogID    string    `json:"blog_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Like      int       `json:"like"`
	Dislike   int       `json:"dislike"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentRequest struct {
	Content string `json:"content" validate:"required"`
}

type CommentResponse struct {
	CommentID string    `json:"comment_id"`
	BlogID    string    `json:"blog_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Like      int       `json:"like"`
	Dislike   int       `json:"dislike"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentUpdateRequest struct {
	Content string `json:"content" validate:"required"`
}

func FromDomainComment(comment *domain.Comment) *CommentJSON {
	return &CommentJSON{
		CommentID: comment.Comment_id,
		BlogID:    comment.Blog_id,
		UserID:    comment.User_id,
		Content:   comment.Content,
		Like:      comment.Like,
		Dislike:   comment.Dislike,
		CreatedAt: comment.Created_at,
		UpdatedAt: comment.Updated_at,
	}
}

func (c *CommentJSON) ToDomain() *domain.Comment {
	return &domain.Comment{
		Comment_id: c.CommentID,
		Blog_id:    c.BlogID,
		User_id:    c.UserID,
		Content:    c.Content,
		Like:       c.Like,
		Dislike:    c.Dislike,
		Created_at: c.CreatedAt,
		Updated_at: c.UpdatedAt,
	}
} 
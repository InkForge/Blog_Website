package models

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type CommentReactionJSON struct {
	CommentID string    `json:"comment_id"`
	UserID    string    `json:"user_id"`
	Action    int       `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentReactionRequest struct {
	Action int `json:"action" validate:"required,oneof=-1 0 1"`
}

type CommentReactionResponse struct {
	CommentID string    `json:"comment_id"`
	UserID    string    `json:"user_id"`
	Action    int       `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromDomainCommentReaction(reaction *domain.CommentReaction) *CommentReactionJSON {
	return &CommentReactionJSON{
		CommentID: reaction.Comment_id,
		UserID:    reaction.User_id,
		Action:    reaction.Action,
		CreatedAt: reaction.Created_at,
		UpdatedAt: reaction.Updated_at,
	}
}

func (c *CommentReactionJSON) ToDomain() *domain.CommentReaction {
	return &domain.CommentReaction{
		Comment_id: c.CommentID,
		User_id:    c.UserID,
		Action:     c.Action,
		Created_at: c.CreatedAt,
		Updated_at: c.UpdatedAt,
	}
} 
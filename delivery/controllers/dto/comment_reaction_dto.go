package dto

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

// ReactionRequest represents the request body for reaction operations
type ReactionRequest struct {
	Action int `json:"action" binding:"required,oneof=-1 0 1"` // -1 for DISLIKE, 0 for NONE, 1 for LIKE
}

// ReactionResponse represents the response for reaction operations
type ReactionResponse struct {
	CommentID string `json:"comment_id"`
	UserID    string `json:"user_id"`
	Action    int    `json:"action"` // -1 for DISLIKE, 0 for NONE, 1 for LIKE
	Message   string `json:"message"`
}

// UserReactionResponse represents the response for getting user reaction
type UserReactionResponse struct {
	CommentID string `json:"comment_id"`
	UserID    string `json:"user_id"`
	Action    int    `json:"action"` // -1 for DISLIKE, 0 for NONE, 1 for LIKE
}

// CommentReactionJSON represents the JSON structure for comment reactions
type CommentReactionJSON struct {
	CommentID string    `json:"comment_id"`
	UserID    string    `json:"user_id"`
	Action    int       `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromDomainCommentReaction converts domain CommentReaction to CommentReactionJSON
func FromDomainCommentReaction(reaction domain.CommentReaction) CommentReactionJSON {
	return CommentReactionJSON{
		CommentID: reaction.Comment_id,
		UserID:    reaction.User_id,
		Action:    reaction.Action,
		CreatedAt: reaction.Created_at,
		UpdatedAt: reaction.Updated_at,
	}
}

// ToDomainCommentReaction converts CommentReactionJSON to domain CommentReaction
func (c *CommentReactionJSON) ToDomain() *domain.CommentReaction {
	return &domain.CommentReaction{
		Comment_id: c.CommentID,
		User_id:    c.UserID,
		Action:     c.Action,
		Created_at: c.CreatedAt,
		Updated_at: c.UpdatedAt,
	}
}

// CreateReactionResponse creates a ReactionResponse with the given parameters
func CreateReactionResponse(commentID, userID string, action int, message string) ReactionResponse {
	return ReactionResponse{
		CommentID: commentID,
		UserID:    userID,
		Action:    action,
		Message:   message,
	}
}

// CreateUserReactionResponse creates a UserReactionResponse with the given parameters
func CreateUserReactionResponse(commentID, userID string, action int) UserReactionResponse {
	return UserReactionResponse{
		CommentID: commentID,
		UserID:    userID,
		Action:    action,
	}
} 
package dto

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

// CommentRequest represents the request body for adding a comment
type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CommentUpdateRequest represents the request body for updating a comment
type CommentUpdateRequest struct {
	Content string `json:"content" binding:"required"`
}

// CommentResponse represents the response for comment operations
type CommentResponse struct {
	ID        string    `json:"id"`
	BlogID    string    `json:"blog_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentListResponse represents the response for listing comments
type CommentListResponse struct {
	Success  bool             `json:"success"`
	Comments []CommentResponse `json:"comments"`
	Count    int              `json:"count"`
}

// FromDomainComment converts domain Comment to CommentResponse
func FromDomainComment(comment domain.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.Comment_id,
		BlogID:    comment.Blog_id,
		UserID:    comment.User_id,
		Content:   comment.Content,
		Likes:     comment.Like,
		Dislikes:  comment.Dislike,
		CreatedAt: comment.Created_at,
		UpdatedAt: comment.Updated_at,
	}
}

// FromDomainComments converts slice of domain Comments to CommentListResponse
func FromDomainComments(comments []domain.Comment) CommentListResponse {
	var responses []CommentResponse
	for _, comment := range comments {
		responses = append(responses, FromDomainComment(comment))
	}

	return CommentListResponse{
		Success:  true,
		Comments: responses,
		Count:    len(responses),
	}
}

// ToDomainComment converts CommentRequest to domain Comment
func (req *CommentRequest) ToDomainComment(blogID, userID string) *domain.Comment {
	return &domain.Comment{
		Blog_id: blogID,
		User_id: userID,
		Content: req.Content,
	}
} 
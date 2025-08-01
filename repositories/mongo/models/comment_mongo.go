package models

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CommentID string             `bson:"comment_id"`
	BlogID    string             `bson:"blog_id"`
	UserID    string             `bson:"user_id"`
	Content   string             `bson:"content"`
	Like      int                `bson:"like"`
	Dislike   int                `bson:"dislike"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func FromDomainComment(comment *domain.Comment) *CommentMongo {
	return &CommentMongo{
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

func (c *CommentMongo) ToDomain() *domain.Comment {
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
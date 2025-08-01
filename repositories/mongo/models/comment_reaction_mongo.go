package models

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentReactionMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CommentID string             `bson:"comment_id"`
	UserID    string             `bson:"user_id"`
	Action    int                `bson:"action"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func FromDomainCommentReaction(reaction *domain.CommentReaction) *CommentReactionMongo {
	return &CommentReactionMongo{
		CommentID: reaction.Comment_id,
		UserID:    reaction.User_id,
		Action:    reaction.Action,
		CreatedAt: reaction.Created_at,
		UpdatedAt: reaction.Updated_at,
	}
}

func (c *CommentReactionMongo) ToDomain() *domain.CommentReaction {
	return &domain.CommentReaction{
		Comment_id: c.CommentID,
		User_id:    c.UserID,
		Action:     c.Action,
		Created_at: c.CreatedAt,
		Updated_at: c.UpdatedAt,
	}
} 
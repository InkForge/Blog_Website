package models

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoBlogReaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Blog_id       string             `bson:"blog_id"`
	User_id       string             `bson:"user_id"`
	Reaction_type int                `bson:"reaction_type"`
	Created_at    time.Time          `bson:"created_at"`
}

func FromDomainBlogReaction(reaction *domain.BlogReaction) (*MongoBlogReaction, error) {
	var objID primitive.ObjectID
	if reaction.ID != "" {
		var err error
		objID, err = primitive.ObjectIDFromHex(reaction.ID)
		if err != nil {
			return nil, domain.ErrInvalidBlogIdFormat
		}
	}

	return &MongoBlogReaction{
		ID:            objID,
		Blog_id:       reaction.Blog_id,
		User_id:       reaction.User_id,
		Reaction_type: reaction.Reaction_type,
		Created_at:    reaction.Created_at,
	}, nil
}

func (mr *MongoBlogReaction) ToDomainBlogReaction() *domain.BlogReaction {
	return &domain.BlogReaction{
		ID:            mr.ID.Hex(),
		Blog_id:       mr.Blog_id,
		User_id:       mr.User_id,
		Reaction_type: mr.Reaction_type,
		Created_at:    mr.Created_at,
	}
}

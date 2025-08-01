package models

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoBlogView struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Blog_id  string             `bson:"blog_id"`
	User_id  string             `bson:"user_id"`
	ViewedAt time.Time          `bson:"viewed_at"`
}

func FromDomainBlogView(view *domain.BlogView) (*MongoBlogView, error) {
	var objID primitive.ObjectID
	if view.ID != "" {
		var err error
		objID, err = primitive.ObjectIDFromHex(view.ID)
		if err != nil {
			return nil, domain.ErrInvalidBlogIdFormat
		}
	}
	return &MongoBlogView{
		ID:       objID,
		Blog_id:  view.Blog_id,
		User_id:  view.User_id,
		ViewedAt: view.ViewedAt,
	}, nil
}

func (mv *MongoBlogView) ToDomainBlogView() *domain.BlogView {
	return &domain.BlogView{
		ID:       mv.ID.Hex(),
		Blog_id:  mv.Blog_id,
		User_id:  mv.User_id,
		ViewedAt: mv.ViewedAt,
	}
}

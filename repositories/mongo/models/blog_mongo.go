package models

import (
	"fmt"
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoBlog struct {
	Blog_id primitive.ObjectID `bson:"_id,omitempty"`
	User_id string             `bson:"user_id"`

	Title   string   `bson:"title"`
	Images  []string `bson:"images"`
	Content string   `bson:"content"`
	Tag_ids []string `bson:"tag_ids"`

	Comment_count int `bson:"comment_count"`
	Like_count    int `bson:"like_count"`
	Dislike_count int `bson:"dislike_count"`
	View_count    int `bson:"view_count"`

	Created_at time.Time `bson:"created_at"`
	Updated_at time.Time `bson:"updated_at"`
}

func FromDomain(blog *domain.Blog) (*MongoBlog, error) {
	var objID primitive.ObjectID
	if blog.Blog_id != "" {
		var err error
		objID, err = primitive.ObjectIDFromHex(blog.Blog_id)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrInvalidBlogIdFormat, err)
		}
	}
	return &MongoBlog{
		Blog_id: objID,
		User_id: blog.User_id,

		Title:   blog.Title,
		Images:  blog.Images,
		Content: blog.Content,
		Tag_ids: blog.Tag_ids,

		Comment_count: blog.Comment_count,
		Like_count:    blog.Like_count,
		Dislike_count: blog.Dislike_count,
		View_count:    blog.View_count,

		Created_at: blog.Created_at,
		Updated_at: blog.Updated_at,
	}, nil
}

func (b *MongoBlog) ToDomain() *domain.Blog {
	return &domain.Blog{

		Blog_id: b.Blog_id.Hex(),
		User_id: b.User_id,

		Title:   b.Title,
		Images:  b.Images,
		Content: b.Content,
		Tag_ids: b.Tag_ids,

		Comment_count: b.Comment_count,
		Like_count:    b.Like_count,
		Dislike_count: b.Dislike_count,
		View_count:    b.View_count,

		Created_at: b.Created_at,
		Updated_at: b.Updated_at,
	}
}

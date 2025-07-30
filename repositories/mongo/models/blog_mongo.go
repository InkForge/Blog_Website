package models

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogMongo struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	BlogID        string             `bson:"blog_id"`
	UserID        string             `bson:"user_id"`
	Title         string             `bson:"title"`
	Images        []string           `bson:"images"`
	Content       string             `bson:"content"`
	TagIDs        []string           `bson:"tag_ids"`
	CommentIDs    []string           `bson:"comment_ids"`
	PostedAt      time.Time          `bson:"posted_at"`
	LikeCounts    int                `bson:"like_counts"`
	DislikeCounts int                `bson:"dislike_counts"`
	ShareCount    int                `bson:"share_count"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

func FromDomain(blog *domain.Blog) *BlogMongo {
	return &BlogMongo{
		BlogID:        blog.Blog_id,
		UserID:        blog.User_id,
		Title:         blog.Title,
		Images:        blog.Images,
		Content:       blog.Content,
		TagIDs:        blog.Tag_ids,
		CommentIDs:    blog.Comment_ids,
		PostedAt:      blog.Posted_at,
		LikeCounts:    blog.Like_counts,
		DislikeCounts: blog.Dislike_counts,
		ShareCount:    blog.Share_count,
		CreatedAt:     blog.Created_at,
		UpdatedAt:     blog.Updated_at,
	}
}

func (b *BlogMongo) ToDomain() *domain.Blog {
	return &domain.Blog{
		Blog_id:        b.BlogID,
		User_id:        b.UserID,
		Title:          b.Title,
		Images:         b.Images,
		Content:        b.Content,
		Tag_ids:        b.TagIDs,
		Comment_ids:    b.CommentIDs,
		Posted_at:      b.PostedAt,
		Like_counts:    b.LikeCounts,
		Dislike_counts: b.DislikeCounts,
		Share_count:    b.ShareCount,
		Created_at:     b.CreatedAt,
		Updated_at:     b.UpdatedAt,
	}
}

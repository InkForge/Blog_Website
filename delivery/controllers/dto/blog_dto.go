package dto

import (
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type PaginationJson struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PaginatedBlogsJson struct {
	Blogs      []BlogJson     `json:"blogs"`
	Pagination PaginationJson `json:"pagination"`
}

type BlogJson struct {
	BlogID       string    `json:"blog_id"`
	UserID       string    `json:"user_id"`
	Title        string    `json:"title"`
	Images       []string  `json:"images"`
	Content      string    `json:"content"`
	TagIDs       []string  `json:"tag_ids"`
	CommentCount int       `json:"comment_count"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	ViewCount    int       `json:"view_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func FromDomainBlog(blog *domain.Blog) *BlogJson {
	return &BlogJson{
		BlogID:       blog.Blog_id,
		UserID:       blog.User_id,
		Title:        blog.Title,
		Images:       blog.Images,
		Content:      blog.Content,
		TagIDs:       blog.Tag_ids,
		CommentCount: blog.Comment_count,
		LikeCount:    blog.Like_count,
		DislikeCount: blog.Dislike_count,
		ViewCount:    blog.View_count,
		CreatedAt:    blog.Created_at,
		UpdatedAt:    blog.Updated_at,
	}
}

func (bj *BlogJson) ToDomainBlog() *domain.Blog {
	return &domain.Blog{
		Blog_id:       bj.BlogID,
		User_id:       bj.UserID,
		Title:         bj.Title,
		Images:        bj.Images,
		Content:       bj.Content,
		Tag_ids:       bj.TagIDs,
		Comment_count: bj.CommentCount,
		Like_count:    bj.LikeCount,
		Dislike_count: bj.DislikeCount,
		View_count:    bj.ViewCount,
		Created_at:    bj.CreatedAt,
		Updated_at:    bj.UpdatedAt,
	}
}

func FromDomainPaginatedBlogs(pb domain.PaginatedBlogs) PaginatedBlogsJson {
	blogs := make([]BlogJson, len(pb.Blogs))
	for i, b := range pb.Blogs {
		blogs[i] = *FromDomainBlog(&b)
	}
	return PaginatedBlogsJson{
		Blogs: blogs,
		Pagination: PaginationJson{
			Page:  pb.Pagination.Page,
			Limit: pb.Pagination.Limit,
			Total: pb.Pagination.Total,
		},
	}
}

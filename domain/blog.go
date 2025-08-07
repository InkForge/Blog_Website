package domain

import (
	"context"
	"time"
)

type Blog struct {
	Blog_id string
	User_id string

	Title   string
	Images  []string
	Content string
	Tag_ids []string

	Comment_count int
	Like_count    int
	Dislike_count int
	View_count    int

	Created_at time.Time
	Updated_at time.Time
}

type Pagination struct {
	Page  int
	Limit int
	Total int
}

type PaginatedBlogs struct {
	Blogs      []Blog
	Pagination Pagination
}

type FilterParams struct {
	TagIDs     []string
	Popularity string
	Page       int
	Limit      int
}

type IBlogRepository interface {
	Create(ctx context.Context, blog Blog) (string, error)
	GetAll(ctx context.Context, page, limit int) ([]Blog, int, error)
	GetByID(ctx context.Context, blogID string) (Blog, error)
	Update(ctx context.Context, blog Blog) error
	Delete(ctx context.Context, blogID string) error

	Search(ctx context.Context, title string, user_ids []string, page, limit int) ([]Blog, int, error)
	Filter(ctx context.Context, params FilterParams) ([]Blog, int, error)

	// Reactions
	IncrementLike(ctx context.Context, blogID string) error
	DecrementLike(ctx context.Context, blogID string) error
	IncrementDisLike(ctx context.Context, blogID string) error
	DecrementDisLike(ctx context.Context, blogID string) error
	ToggleLikeDislikeCounts(ctx context.Context, blogID string, to_like, to_dislike int) error

	// Views
	IncrementView(ctx context.Context, blogID string) error

	// Comments
	AddCommentID(ctx context.Context, blogID, commentID string) error
	RemoveCommentID(ctx context.Context, blogID, commentID string) error
}

type IBlogUseCase interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetAllBlogs(ctx context.Context, page, limit int) (*PaginatedBlogs, error)
	GetBlogByID(ctx context.Context, blogID, userID string) (Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog, userID string) error
	DeleteBlog(ctx context.Context, blogID string) error

	SearchBlogs(ctx context.Context, title, author string, page, limit int) (*PaginatedBlogs, error)
	FilterBlogs(ctx context.Context, params FilterParams) (*PaginatedBlogs, error)
}

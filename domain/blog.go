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

type IBlogRepository interface {
	Create(ctx context.Context, blog Blog) (string, error)
	GetAll(ctx context.Context) ([]Blog, error)
	GetByID(ctx context.Context, blogID string) (Blog, error)
	Update(ctx context.Context, blog Blog) error
	Delete(ctx context.Context, blogID string) error
	// receives matching (first names or last names user_ids) and title
	Search(ctx context.Context, title string, user_ids []string) ([]Blog, error)

	// Operations related to blog_reaction
	IncrementLike(ctx context.Context, blogID string) error
	DecrementLike(ctx context.Context, blogID string) error
	IncrementDisLike(ctx context.Context, blogID string) error
	DecrementDisLike(ctx context.Context, blogID string) error
	ToggleLikeDislikeCounts(ctx context.Context, blogID string, to_like, to_dislike int) error
	
	// Operations related to blog views
	IncrementView(ctx context.Context, blogID string) error
	
	// Operations related to comments
	AddCommentID(ctx context.Context, blogID, commentID string) error
	RemoveCommentID(ctx context.Context, blogID, commentID string) error
}

type IBlogUsecase interface {
	CreateBlog(ctx context.Context, blog *Blog) (string, error)
	GetAllBlogs(ctx context.Context) ([]Blog, error)
	GetBlogByID(ctx context.Context, blogID, userID string) (*Blog, error)
	UpdateBlog(ctx context.Context, blog *Blog) error
	DeleteBlog(ctx context.Context, blogID string) error
}

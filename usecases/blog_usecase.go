package usecases

import (
	"context"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type BlogUsecase struct {
	blogRepository domain.IBlogRepository
	contextTimeout time.Duration
}

func NewBlogUsecase(blogRepo domain.IBlogRepository, timeout time.Duration) domain.IBlogUsecase {
	return &BlogUsecase{
		blogRepository: blogRepo,
		contextTimeout: timeout,
	}
}

func (bu *BlogUsecase) CreateBlog(ctx context.Context, blog *domain.Blog) (string, error) {
	if blog == nil {
		return "", domain.ErrBlogRequired
	}
	if blog.Title == "" {
		return "", domain.ErrEmptyTitle
	}
	if blog.Content == "" {
		return "", domain.ErrEmptyContent
	}
	if blog.User_id == "" {
		return "", domain.ErrInvalidUserID
	}

	blog.Created_at = time.Now()
	blog.Updated_at = blog.Created_at
	blog.Posted_at = blog.Created_at
	blog.Like_counts = 0
	blog.Dislike_counts = 0
	blog.Share_count = 0

	ctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()

	return bu.blogRepository.Create(ctx, *blog)
}

func (bu *BlogUsecase) GetAllBlogs(ctx context.Context) ([]domain.Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()

	return bu.blogRepository.GetAll(ctx)
}

func (bu *BlogUsecase) GetBlogByID(ctx context.Context, blogID string) (*domain.Blog, error) {
	if blogID == "" {
		return nil, domain.ErrInvalidBlogID
	}

	ctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()

	blog, err := bu.blogRepository.GetByID(ctx, blogID)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (bu *BlogUsecase) UpdateBlog(ctx context.Context, blog *domain.Blog) error {
	if blog == nil {
		return domain.ErrBlogRequired
	}
	if blog.Blog_id == "" {
		return domain.ErrBlogIDRequired
	}

	ctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()

	// Get existing blog to preserve unchanged fields
	existing, err := bu.blogRepository.GetByID(ctx, blog.Blog_id)
	if err != nil {
		return err
	}

	// Track if any changes were made
	changesMade := false

	// Apply partial updates - only update user-editable fields
	if blog.Title != "" {
		existing.Title = blog.Title
		changesMade = true
	}
	if blog.Content != "" {
		existing.Content = blog.Content
		changesMade = true
	}
	if blog.User_id != "" {
		existing.User_id = blog.User_id
		changesMade = true
	}
	if blog.Images != nil {
		existing.Images = blog.Images
		changesMade = true
	}
	if blog.Tag_ids != nil {
		existing.Tag_ids = blog.Tag_ids
		changesMade = true
	}
	if !blog.Posted_at.IsZero() {
		existing.Posted_at = blog.Posted_at
		changesMade = true
	}

	// Check if any changes were actually made
	if !changesMade {
		return domain.ErrNoBlogChangesMade
	}

	// Update timestamp
	existing.Updated_at = time.Now()

	return bu.blogRepository.Update(ctx, existing)
}

func (bu *BlogUsecase) DeleteBlog(ctx context.Context, blogID string) error {
	if blogID == "" {
		return domain.ErrInvalidBlogID
	}

	ctx, cancel := context.WithTimeout(ctx, bu.contextTimeout)
	defer cancel()

	return bu.blogRepository.Delete(ctx, blogID)
}

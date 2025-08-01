package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type BlogUsecase struct {
	blogRepo           domain.IBlogRepository
	blogViewRepo       domain.IBlogViewRepository
	transactionManager domain.TransactionManager
}

func NewBlogUsecase(blogRepo domain.IBlogRepository, blogViewRepo domain.IBlogViewRepository) *BlogUsecase {
	return &BlogUsecase{
		blogRepo:     blogRepo,
		blogViewRepo: blogViewRepo,
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
	blog.Like_count = 0
	blog.Dislike_count = 0
	blog.View_count = 0

	return bu.blogRepo.Create(ctx, *blog)
}

func (bu *BlogUsecase) GetAllBlogs(ctx context.Context) ([]domain.Blog, error) {
	return bu.blogRepo.GetAll(ctx)
}

func (bu *BlogUsecase) GetBlogByID(ctx context.Context, blogID, userID string) (*domain.Blog, error) {
	if blogID == "" {
		return nil, domain.ErrBlogIDRequired
	}
	if userID == "" {
		return nil, domain.ErrInvalidUserID
	}

	var fetchedBlog domain.Blog
	// All operations for retrieving and logging a view are wrapped in a single transaction
	err := bu.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		blog, err := bu.blogRepo.GetByID(txCtx, blogID)
		if err != nil {
			return err
		}
		fetchedBlog = blog

		err = bu.blogViewRepo.CreateViewRecord(txCtx, blogID, userID)
		if err != nil {
			if errors.Is(err, domain.ErrViewRecordAlreadyExists) {
				return domain.ErrViewRecordAlreadyExists
			}
			return domain.ErrCreateViewRecordFailed
		}

		err = bu.blogRepo.IncrementView(txCtx, blogID)
		if err != nil {
			return domain.ErrIncrementViewFailed
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &fetchedBlog, nil
}

func (bu *BlogUsecase) UpdateBlog(ctx context.Context, blog *domain.Blog) error {
	if blog == nil {
		return domain.ErrBlogRequired
	}
	if blog.Blog_id == "" {
		return domain.ErrBlogIDRequired
	}

	existing, err := bu.blogRepo.GetByID(ctx, blog.Blog_id)
	if err != nil {
		return err
	}

	if blog.Title != "" {
		existing.Title = blog.Title
	}
	if blog.Content != "" {
		existing.Content = blog.Content
	}
	if blog.User_id != "" {
		existing.User_id = blog.User_id
	}
	if blog.Images != nil {
		existing.Images = blog.Images
	}
	if blog.Tag_ids != nil {
		existing.Tag_ids = blog.Tag_ids
	}

	existing.Updated_at = time.Now()

	return bu.blogRepo.Update(ctx, existing)
}

func (bu *BlogUsecase) DeleteBlog(ctx context.Context, blogID string) error {
	if blogID == "" {
		return domain.ErrBlogIDRequired
	}
	return bu.blogRepo.Delete(ctx, blogID)
}

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
	tagRepo            domain.ITagRepository
	userRepo           domain.IUserRepository
	transactionManager domain.ITransactionManager
}

func NewBlogUsecase(blogRepo domain.IBlogRepository,
	blogViewRepo domain.IBlogViewRepository,
	tagRepo domain.ITagRepository,
	userRepo domain.IUserRepository,
	transactionManager domain.ITransactionManager,
) domain.IBlogUseCase {
	return &BlogUsecase{
		blogRepo:           blogRepo,
		blogViewRepo:       blogViewRepo,
		tagRepo:            tagRepo,
		userRepo:           userRepo,
		transactionManager: transactionManager,
	}
}

func (bu *BlogUsecase) ensureTagsExist(ctx context.Context, tagNames []string) ([]string, error) {
	if len(tagNames) == 0 {
		return nil, nil
	}

	// Find existing tags
	existingTags, err := bu.tagRepo.FindByNames(ctx, tagNames)
	if err != nil {
		return nil, err
	}
	existingTagNames := make(map[string]bool)
	existingTagIDs := make([]string, 0, len(existingTags))
	for _, t := range existingTags {
		existingTagNames[t.TagName] = true
		existingTagIDs = append(existingTagIDs, t.Tag_id)
	}

	// Filter out names that already exist to find new tags to create
	var toCreate []string
	for _, name := range tagNames {
		if !existingTagNames[name] {
			toCreate = append(toCreate, name)
		}
	}

	// Create only new tags
	var createdTags []domain.Tag
	if len(toCreate) > 0 {
		createdTags, err = bu.tagRepo.CreateMany(ctx, toCreate)
		if err != nil {
			return nil, err
		}
	}

	allTagIDs := make([]string, 0, len(existingTags)+len(createdTags))
	allTagIDs = append(allTagIDs, existingTagIDs...)
	for _, t := range createdTags {
		allTagIDs = append(allTagIDs, t.Tag_id)
	}

	return allTagIDs, nil
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

	var blogID string
	err := bu.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if len(blog.Tag_ids) > 0 {
			allTagIDs, err := bu.ensureTagsExist(txCtx, blog.Tag_ids)
			if err != nil {
				return err
			}
			blog.Tag_ids = allTagIDs
		}

		now := time.Now()
		blog.Created_at = now
		blog.Updated_at = now
		blog.Like_count = 0
		blog.Dislike_count = 0
		blog.View_count = 0
		blog.Comment_count = 0

		var err error
		blogID, err = bu.blogRepo.Create(txCtx, *blog)
		if err != nil {
			return err
		}

		return nil
	})

	return blogID, err
}

func (bu *BlogUsecase) GetAllBlogs(ctx context.Context, page, limit int) (*domain.PaginatedBlogs, error) {
	blogs, total, err := bu.blogRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	paginatedBlogs := &domain.PaginatedBlogs{
		Blogs: blogs,
		Pagination: domain.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}

	return paginatedBlogs, nil
}

func (bu *BlogUsecase) GetBlogByID(ctx context.Context, blogID, userID string) (*domain.Blog, error) {
	var fetchedBlog domain.Blog

	err := bu.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		blog, err := bu.blogRepo.GetByID(txCtx, blogID)
		if err != nil {
			return err
		}
		fetchedBlog = blog

		err = bu.blogViewRepo.CreateViewRecord(txCtx, blogID, userID)
		if err != nil {
			if errors.Is(err, domain.ErrViewRecordAlreadyExists) {
				return nil
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

	return bu.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		existing, err := bu.blogRepo.GetByID(txCtx, blog.Blog_id)
		if err != nil {
			return err
		}

		if blog.Title != "" {
			existing.Title = blog.Title
		}
		if blog.Content != "" {
			existing.Content = blog.Content
		}
		if blog.Images != nil {
			existing.Images = blog.Images
		}

		if blog.Tag_ids != nil {
			allTagIDs, err := bu.ensureTagsExist(txCtx, blog.Tag_ids)
			if err != nil {
				return err
			}
			existing.Tag_ids = allTagIDs
		}

		existing.Updated_at = time.Now()

		return bu.blogRepo.Update(txCtx, existing)
	})
}

func (bu *BlogUsecase) DeleteBlog(ctx context.Context, blogID string) error {
	if blogID == "" {
		return domain.ErrBlogIDRequired
	}
	return bu.blogRepo.Delete(ctx, blogID)

}

func (bu *BlogUsecase) FilterBlogs(ctx context.Context, params domain.FilterParams) (*domain.PaginatedBlogs, error) {
	blogs, total, err := bu.blogRepo.Filter(ctx, params)
	if err != nil {
		return nil, err
	}
	paginatedBlogs := &domain.PaginatedBlogs{
		Blogs: blogs,
		Pagination: domain.Pagination{
			Page:  params.Page,
			Limit: params.Limit,
			Total: total,
		},
	}
	return paginatedBlogs, nil
}

func (bu *BlogUsecase) SearchBlogs(ctx context.Context, title, author string, page, limit int) (*domain.PaginatedBlogs, error) {
	var userIDs []string
	if author != "" && bu.userRepo != nil {
		users, err := bu.userRepo.FindUsersByName(ctx, author)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			userIDs = append(userIDs, u.UserID)
		}
	}
	blogs, total, err := bu.blogRepo.Search(ctx, title, userIDs, page, limit)
	if err != nil {
		return nil, err
	}
	paginatedBlogs := &domain.PaginatedBlogs{
		Blogs: blogs,
		Pagination: domain.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
	return paginatedBlogs, nil
}

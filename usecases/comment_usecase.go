package usecases

import (
	"context"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type CommentUsecase struct {
	blogRepository    domain.IBlogRepository
	commentRepository domain.ICommentRepository
	contextTimeout    time.Duration
}

func NewCommentUsecase(blogRepo domain.IBlogRepository, commentRepo domain.ICommentRepository, timeout time.Duration) *CommentUsecase {
	return &CommentUsecase{
		blogRepository:    blogRepo,
		commentRepository: commentRepo,
		contextTimeout:    timeout,
	}
}

func (cu *CommentUsecase) AddComment(ctx context.Context, blogID string, comment *domain.Comment) error {
	if comment == nil {
		return domain.ErrCommentRequired
	}
	if blogID == "" {
		return domain.ErrInvalidBlogID
	}
	if comment.Content == "" {
		return domain.ErrEmptyCommentContent
	}
	if comment.User_id == "" {
		return domain.ErrInvalidUserID
	}

	// Verify blog exists
	_, err := cu.blogRepository.GetByID(ctx, blogID)
	if err != nil {
		return domain.ErrBlogNotFound
	}

	comment.Blog_id = blogID
	comment.Created_at = time.Now()
	comment.Updated_at = comment.Created_at

	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	commentID, err := cu.commentRepository.Create(ctx, *comment)
	if err != nil {
		return err
	}

	// Update blog's comment_ids array
	blog, err := cu.blogRepository.GetByID(ctx, blogID)
	if err != nil {
		return err
	}

	blog.Comment_ids = append(blog.Comment_ids, commentID)
	blog.Updated_at = time.Now()

	return cu.blogRepository.Update(ctx, blog)
}

func (cu *CommentUsecase) RemoveComment(ctx context.Context, blogID, commentID string) error {
	if blogID == "" {
		return domain.ErrInvalidBlogID
	}
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}

	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	// Verify comment exists and belongs to the blog
	comment, err := cu.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return domain.ErrCommentNotFound
	}
	if comment.Blog_id != blogID {
		return domain.ErrCommentNotFound
	}

	// Delete the comment
	err = cu.commentRepository.Delete(ctx, commentID)
	if err != nil {
		return err
	}

	// Update blog's comment_ids array
	blog, err := cu.blogRepository.GetByID(ctx, blogID)
	if err != nil {
		return err
	}

	// Remove commentID from the array
	var newCommentIDs []string
	for _, id := range blog.Comment_ids {
		if id != commentID {
			newCommentIDs = append(newCommentIDs, id)
		}
	}
	blog.Comment_ids = newCommentIDs
	blog.Updated_at = time.Now()

	return cu.blogRepository.Update(ctx, blog)
}

func (cu *CommentUsecase) GetBlogComments(ctx context.Context, blogID string) ([]domain.Comment, error) {
	if blogID == "" {
		return nil, domain.ErrInvalidBlogID
	}

	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	// Verify blog exists
	_, err := cu.blogRepository.GetByID(ctx, blogID)
	if err != nil {
		return nil, domain.ErrBlogNotFound
	}

	return cu.commentRepository.GetByBlogID(ctx, blogID)
}

func (cu *CommentUsecase) UpdateComment(ctx context.Context, commentID string, comment *domain.Comment) error {
	if comment == nil {
		return domain.ErrCommentRequired
	}
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}
	if comment.Content == "" {
		return domain.ErrEmptyCommentContent
	}

	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	// Get existing comment
	existing, err := cu.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return domain.ErrCommentNotFound
	}

	// Update content
	existing.Content = comment.Content
	existing.Updated_at = time.Now()

	return cu.commentRepository.Update(ctx, existing)
} 
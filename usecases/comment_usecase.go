package usecases

import (
	"context"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type CommentUsecase struct {
	blogRepository       domain.IBlogRepository
	commentRepository    domain.ICommentRepository
	transactionManager   domain.ITransactionManager
	contextTimeout       time.Duration
}

func NewCommentUsecase(blogRepo domain.IBlogRepository, commentRepo domain.ICommentRepository, transactionManager domain.ITransactionManager, timeout time.Duration) domain.ICommentUsecase {
	return &CommentUsecase{
		blogRepository:     blogRepo,
		commentRepository:  commentRepo,
		transactionManager: transactionManager,
		contextTimeout:     timeout,
	}
}

func (cu *CommentUsecase) AddComment(ctx context.Context, blogID string, comment *domain.Comment) (string, error) {
	if comment == nil {
		return "", domain.ErrCommentRequired
	}
	if blogID == "" {
		return "", domain.ErrInvalidBlogID
	}
	if comment.Content == "" {
		return "", domain.ErrEmptyCommentContent
	}
	if comment.User_id == "" {
		return "", domain.ErrInvalidUserID
	}

	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	// ensure blog exists
	_, err := cu.blogRepository.GetByID(ctx, blogID)
	if err != nil {
		return "", domain.ErrBlogNotFound
	}

	comment.Blog_id = blogID
	comment.Like = 0
	comment.Dislike = 0
	comment.Created_at = time.Now()
	comment.Updated_at = comment.Created_at

	var commentID string
	err = cu.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Create the comment
		id, err := cu.commentRepository.Create(txCtx, *comment)
		if err != nil {
			return err
		}
		commentID = id

		// Add comment ID to blog
		err = cu.blogRepository.AddCommentID(txCtx, blogID, commentID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return commentID, nil
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

	comment, err := cu.commentRepository.GetByID(ctx, commentID)
	if err != nil || comment.Blog_id != blogID {
		return domain.ErrCommentNotFound
	}

	err = cu.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Delete the comment
		err := cu.commentRepository.Delete(txCtx, commentID)
		if err != nil {
			return err
		}

		// Remove comment ID from blog
		err = cu.blogRepository.RemoveCommentID(txCtx, blogID, commentID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (cu *CommentUsecase) GetBlogComments(ctx context.Context, blogID string) ([]domain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	_, err := cu.blogRepository.GetByID(ctx, blogID)
	if err != nil {
		return nil, domain.ErrBlogNotFound
	}

	return cu.commentRepository.GetByBlogID(ctx, blogID)
}

func (cu *CommentUsecase) UpdateComment(ctx context.Context, commentID string, comment *domain.Comment) error {
	if comment == nil || commentID == "" || comment.Content == "" {
		return domain.ErrCommentRequired
	}

	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	existing, err := cu.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return domain.ErrCommentNotFound
	}

	existing.Content = comment.Content
	existing.Updated_at = time.Now()

	return cu.commentRepository.Update(ctx, existing)
} 
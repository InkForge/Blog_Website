package usecases

import (
	"context"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type CommentUsecase struct {
	blogRepository     domain.IBlogRepository
	commentRepository  domain.ICommentRepository
	transactionManager domain.ITransactionManager
}

func NewCommentUsecase(
	blogRepo domain.IBlogRepository,
	commentRepo domain.ICommentRepository,
	txManager domain.ITransactionManager,
) domain.ICommentUsecase {
	return &CommentUsecase{
		blogRepository:     blogRepo,
		commentRepository:  commentRepo,
		transactionManager: txManager,
	}
}

func (cu *CommentUsecase) AddComment(
	ctx context.Context,
	blogID string,
	comment *domain.Comment,
	role string,
) (string, error) {
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
		id, err := cu.commentRepository.Create(txCtx, *comment)
		if err != nil {
			return err
		}
		commentID = id

		return cu.blogRepository.AddCommentID(txCtx, blogID, commentID)
	})
	if err != nil {
		return "", err
	}

	return commentID, nil
}

func (cu *CommentUsecase) RemoveComment(
	ctx context.Context,
	blogID, commentID, requesterID, role string,
) error {
	if blogID == "" {
		return domain.ErrInvalidBlogID
	}
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}

	comment, err := cu.commentRepository.GetByID(ctx, commentID)
	if err != nil || comment.Blog_id != blogID {
		return domain.ErrCommentNotFound
	}

	if role != "admin" && comment.User_id != requesterID {
		return domain.ErrForbidden
	}

	return cu.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := cu.commentRepository.Delete(txCtx, commentID); err != nil {
			return err
		}
		return cu.blogRepository.RemoveCommentID(txCtx, blogID, commentID)
	})
}

func (cu *CommentUsecase) GetBlogComments(
	ctx context.Context,
	blogID string,
) ([]domain.Comment, error) {
	_, err := cu.blogRepository.GetByID(ctx, blogID)
	if err != nil {
		return nil, domain.ErrBlogNotFound
	}
	return cu.commentRepository.GetByBlogID(ctx, blogID)
}

func (cu *CommentUsecase) UpdateComment(
	ctx context.Context,
	commentID string,
	comment *domain.Comment,
	role string,
) error {
	if comment == nil || commentID == "" || comment.Content == "" {
		return domain.ErrCommentRequired
	}

	existing, err := cu.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return domain.ErrCommentNotFound
	}

	if role != "admin" && comment.User_id != existing.User_id {
		return domain.ErrForbidden
	}

	existing.Content = comment.Content
	existing.Updated_at = time.Now()

	return cu.commentRepository.Update(ctx, existing)
}

func (cu *CommentUsecase) GetCommentByID(
	ctx context.Context,
	commentID string,
) (domain.Comment, error) {
	return cu.commentRepository.GetByID(ctx, commentID)
}

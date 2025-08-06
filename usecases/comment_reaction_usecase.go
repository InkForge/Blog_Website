package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type CommentReactionUsecase struct {
	commentRepository         domain.ICommentRepository
	commentReactionRepository domain.ICommentReactionRepository
	transactionManager        domain.ITransactionManager
}

// NewCommentReactionUsecase creates a new comment reaction use case instance with required dependencies
func NewCommentReactionUsecase(
	commentRepo domain.ICommentRepository,
	reactionRepo domain.ICommentReactionRepository,
	transactionManager domain.ITransactionManager,
) domain.ICommentReactionUsecase {
	return &CommentReactionUsecase{
		commentRepository:         commentRepo,
		commentReactionRepository: reactionRepo,
		transactionManager:        transactionManager,
	}
}

// LikeComment handles liking a comment with transaction support to ensure reaction count consistency
func (cru *CommentReactionUsecase) LikeComment(ctx context.Context, commentID, userID string) error {
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	if _, err := cru.commentRepository.GetByID(ctx, commentID); err != nil {
		if errors.Is(err, domain.ErrCommentNotFound) {
			return err
		}
		return err
	}

	existing, err := cru.commentReactionRepository.GetByCommentAndUser(ctx, commentID, userID)
	if err != nil && !errors.Is(err, domain.ErrCommentReactionNotFound) {
		return err
	}

	now := time.Now()

	return cru.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if errors.Is(err, domain.ErrCommentReactionNotFound) {
			// No existing reaction, create one
			newReaction := domain.CommentReaction{
				Comment_id: commentID,
				User_id:    userID,
				Action:     1,
				Created_at: now,
				Updated_at: now,
			}
			if err := cru.commentReactionRepository.Create(txCtx, newReaction); err != nil {
				return err
			}
		} else {
			if existing.Action == 1 {
				// Already liked, remove it
				if err := cru.commentReactionRepository.Delete(txCtx, commentID, userID); err != nil {
					return err
				}
			} else {
				// Switch to like
				existing.Action = 1
				existing.Updated_at = now
				if err := cru.commentReactionRepository.Update(txCtx, existing); err != nil {
					return err
				}
			}
		}

		return cru.updateCommentReactionCounts(txCtx, commentID)
	})
}

// DislikeComment handles disliking a comment with transaction support to ensure reaction count consistency
func (cru *CommentReactionUsecase) DislikeComment(ctx context.Context, commentID, userID string) error {
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	if _, err := cru.commentRepository.GetByID(ctx, commentID); err != nil {
		if errors.Is(err, domain.ErrCommentNotFound) {
			return err
		}
		return err
	}

	existing, err := cru.commentReactionRepository.GetByCommentAndUser(ctx, commentID, userID)
	if err != nil && !errors.Is(err, domain.ErrCommentReactionNotFound) {
		return err
	}

	now := time.Now()

	return cru.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if errors.Is(err, domain.ErrCommentReactionNotFound) {
			newReaction := domain.CommentReaction{
				Comment_id: commentID,
				User_id:    userID,
				Action:     -1,
				Created_at: now,
				Updated_at: now,
			}
			if err := cru.commentReactionRepository.Create(txCtx, newReaction); err != nil {
				return err
			}
		} else {
			if existing.Action == -1 {
				if err := cru.commentReactionRepository.Delete(txCtx, commentID, userID); err != nil {
					return err
				}
			} else {
				existing.Action = -1
				existing.Updated_at = now
				if err := cru.commentReactionRepository.Update(txCtx, existing); err != nil {
					return err
				}
			}
		}

		return cru.updateCommentReactionCounts(txCtx, commentID)
	})
}

// RemoveReaction removes a user's reaction from a comment with transaction support
func (cru *CommentReactionUsecase) RemoveReaction(ctx context.Context, commentID, userID, role string) error {
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	if _, err := cru.commentRepository.GetByID(ctx, commentID); err != nil {
		if errors.Is(err, domain.ErrCommentNotFound) {
			return err
		}
		return err
	}

	// Check if user has permission to remove this reaction
	existingReaction, err := cru.commentReactionRepository.GetByCommentAndUser(ctx, commentID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCommentReactionNotFound) {
			return domain.ErrCommentReactionNotFound
		}
		return err
	}

	// Only the user who created the reaction or an admin can remove it
	if role != "admin" && existingReaction.User_id != userID {
		return domain.ErrForbidden
	}

	return cru.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := cru.commentReactionRepository.Delete(txCtx, commentID, userID); err != nil {
			return err
		}

		return cru.updateCommentReactionCounts(txCtx, commentID)
	})
}

// GetUserReaction retrieves a user's reaction to a specific comment
func (cru *CommentReactionUsecase) GetUserReaction(ctx context.Context, commentID, userID string) (int, error) {
	if commentID == "" {
		return 0, domain.ErrInvalidCommentID
	}
	if userID == "" {
		return 0, domain.ErrInvalidUserID
	}

	if _, err := cru.commentRepository.GetByID(ctx, commentID); err != nil {
		if errors.Is(err, domain.ErrCommentNotFound) {
			return 0, err
		}
		return 0, err
	}

	reaction, err := cru.commentReactionRepository.GetByCommentAndUser(ctx, commentID, userID)
	if errors.Is(err, domain.ErrCommentReactionNotFound) {
		return 0, nil // No reaction
	}
	if err != nil {
		return 0, err
	}

	return reaction.Action, nil
}

// updateCommentReactionCounts updates the like and dislike counts for a comment based on current reactions
func (cru *CommentReactionUsecase) updateCommentReactionCounts(ctx context.Context, commentID string) error {
	likes, dislikes, err := cru.commentReactionRepository.GetReactionCounts(ctx, commentID)
	if err != nil {
		return err
	}
	return cru.commentRepository.UpdateReactionCounts(ctx, commentID, likes, dislikes)
}

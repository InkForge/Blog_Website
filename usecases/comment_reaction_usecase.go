package usecases

import (
	"context"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type CommentReactionUsecase struct {
	commentRepository        domain.ICommentRepository
	commentReactionRepository domain.ICommentReactionRepository
	contextTimeout           time.Duration
}

func NewCommentReactionUsecase(commentRepo domain.ICommentRepository, reactionRepo domain.ICommentReactionRepository, timeout time.Duration) domain.ICommentReactionUsecase {
	return &CommentReactionUsecase{
		commentRepository:        commentRepo,
		commentReactionRepository: reactionRepo,
		contextTimeout:           timeout,
	}
}

func (cru *CommentReactionUsecase) LikeComment(ctx context.Context, commentID, userID string) error {
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	ctx, cancel := context.WithTimeout(ctx, cru.contextTimeout)
	defer cancel()

	// Verify comment exists
	_, err := cru.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return domain.ErrCommentNotFound
	}

	// Check if user already has a reaction
	existingReaction, err := cru.commentReactionRepository.GetByCommentAndUser(ctx, commentID, userID)
	if err != nil && err != domain.ErrCommentReactionNotFound {
		return err
	}

	now := time.Now()

	if err == domain.ErrCommentReactionNotFound {
		// Create new like reaction
		reaction := domain.CommentReaction{
			Comment_id: commentID,
			User_id:    userID,
			Action:     1, // LIKE
			Created_at: now,
			Updated_at: now,
		}
		err = cru.commentReactionRepository.Create(ctx, reaction)
		if err != nil {
			return err
		}
	} else {
		// Update existing reaction
		if existingReaction.Action == 1 {
			// Already liked, remove the like
			err = cru.commentReactionRepository.Delete(ctx, commentID, userID)
			if err != nil {
				return err
			}
		} else {
			// Change to like
			existingReaction.Action = 1
			existingReaction.Updated_at = now
			err = cru.commentReactionRepository.Update(ctx, existingReaction)
			if err != nil {
				return err
			}
		}
	}

	// Update comment reaction counts
	return cru.updateCommentReactionCounts(ctx, commentID)
}

func (cru *CommentReactionUsecase) DislikeComment(ctx context.Context, commentID, userID string) error {
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	ctx, cancel := context.WithTimeout(ctx, cru.contextTimeout)
	defer cancel()

	// Verify comment exists
	_, err := cru.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return domain.ErrCommentNotFound
	}

	// Check if user already has a reaction
	existingReaction, err := cru.commentReactionRepository.GetByCommentAndUser(ctx, commentID, userID)
	if err != nil && err != domain.ErrCommentReactionNotFound {
		return err
	}

	now := time.Now()

	if err == domain.ErrCommentReactionNotFound {
		// Create new dislike reaction
		reaction := domain.CommentReaction{
			Comment_id: commentID,
			User_id:    userID,
			Action:     -1, // DISLIKE
			Created_at: now,
			Updated_at: now,
		}
		err = cru.commentReactionRepository.Create(ctx, reaction)
		if err != nil {
			return err
		}
	} else {
		// Update existing reaction
		if existingReaction.Action == -1 {
			// Already disliked, remove the dislike
			err = cru.commentReactionRepository.Delete(ctx, commentID, userID)
			if err != nil {
				return err
			}
		} else {
			// Change to dislike
			existingReaction.Action = -1
			existingReaction.Updated_at = now
			err = cru.commentReactionRepository.Update(ctx, existingReaction)
			if err != nil {
				return err
			}
		}
	}

	// Update comment reaction counts
	return cru.updateCommentReactionCounts(ctx, commentID)
}

func (cru *CommentReactionUsecase) RemoveReaction(ctx context.Context, commentID, userID string) error {
	if commentID == "" {
		return domain.ErrInvalidCommentID
	}
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	ctx, cancel := context.WithTimeout(ctx, cru.contextTimeout)
	defer cancel()

	// Verify comment exists
	_, err := cru.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return domain.ErrCommentNotFound
	}

	// Remove the reaction
	err = cru.commentReactionRepository.Delete(ctx, commentID, userID)
	if err != nil {
		return err
	}

	// Update comment reaction counts
	return cru.updateCommentReactionCounts(ctx, commentID)
}

func (cru *CommentReactionUsecase) GetUserReaction(ctx context.Context, commentID, userID string) (int, error) {
	if commentID == "" {
		return 0, domain.ErrInvalidCommentID
	}
	if userID == "" {
		return 0, domain.ErrInvalidUserID
	}

	ctx, cancel := context.WithTimeout(ctx, cru.contextTimeout)
	defer cancel()

	// Verify comment exists
	_, err := cru.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return 0, domain.ErrCommentNotFound
	}

	reaction, err := cru.commentReactionRepository.GetByCommentAndUser(ctx, commentID, userID)
	if err != nil {
		if err == domain.ErrCommentReactionNotFound {
			return 0, nil // No reaction
		}
		return 0, err
	}

	return reaction.Action, nil
}

func (cru *CommentReactionUsecase) updateCommentReactionCounts(ctx context.Context, commentID string) error {
	likeCount, dislikeCount, err := cru.commentReactionRepository.GetReactionCounts(ctx, commentID)
	if err != nil {
		return err
	}

	return cru.commentRepository.UpdateReactionCounts(ctx, commentID, likeCount, dislikeCount)
} 
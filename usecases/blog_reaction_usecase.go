package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type BlogReactionUseCase struct {
	blogRepo           domain.IBlogRepository
	blogReactionRepo   domain.IBlogReactionRepository
	transactionManager domain.TransactionManager
}

func NewBlogReactionUseCase(
	blogRepo domain.IBlogRepository,
	blogReactionRepo domain.IBlogReactionRepository,
	transactionManager domain.TransactionManager,
) *BlogReactionUseCase {
	return &BlogReactionUseCase{
		blogRepo:           blogRepo,
		blogReactionRepo:   blogReactionRepo,
		transactionManager: transactionManager,
	}
}

func (uc *BlogReactionUseCase) LikeBlog(ctx context.Context, blog_id, user_id string) error {

	return uc.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// checking if we are creating a new record or updating existing
		existingReaction, err := uc.blogReactionRepo.GetReactionByUserAndBlog(txCtx, blog_id, user_id)

		if err != nil && !errors.Is(err, domain.ErrBlogReactionNotFound) {
			return domain.ErrCheckBlogReactionFailed
		}

		if errors.Is(err, domain.ErrBlogReactionNotFound) {
			// if no record exists create one
			newReaction := domain.BlogReaction{
				Blog_id:       blog_id,
				User_id:       user_id,
				Reaction_type: 1,
				Created_at:    time.Now(),
			}
			if err := uc.blogReactionRepo.CreateReaction(txCtx, newReaction); err != nil {
				return domain.ErrCreateBlogReactionFailed
			}
			if err := uc.blogRepo.IncrementLike(txCtx, blog_id); err != nil {
				return domain.ErrIncrementLikeFailed
			}
		} else {
			// else update dislike to like and change the count
			if existingReaction.Reaction_type == 1 {
				return nil
			} else {
				existingReaction.Reaction_type = 1
				if err := uc.blogReactionRepo.UpdateReaction(txCtx, existingReaction); err != nil {
					return domain.ErrUpdateBlogReactionFailed
				}
				// +1 like_count and -1 dislike_count
				if err := uc.blogRepo.ToggleLikeDislikeCounts(txCtx, blog_id, 1, -1); err != nil {
					return domain.ErrToggleLikeDislikeFailed
				}
			}
		}
		return nil
	})
}

func (uc *BlogReactionUseCase) DisLikeBlog(ctx context.Context, blog_id, user_id string) error {

	return uc.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// checking if we are creating a new record or updating existing
		existingReaction, err := uc.blogReactionRepo.GetReactionByUserAndBlog(txCtx, blog_id, user_id)
		if err != nil && !errors.Is(err, domain.ErrBlogReactionNotFound) {
			return domain.ErrCheckBlogReactionFailed
		}

		if errors.Is(err, domain.ErrBlogReactionNotFound) {
			// if no record exists create one
			newReaction := domain.BlogReaction{
				Blog_id:       blog_id,
				User_id:       user_id,
				Reaction_type: -1,
				Created_at:    time.Now(),
			}
			if err := uc.blogReactionRepo.CreateReaction(txCtx, newReaction); err != nil {
				return domain.ErrCreateBlogReactionFailed
			}
			if err := uc.blogRepo.IncrementDisLike(txCtx, blog_id); err != nil {
				return domain.ErrIncrementLikeFailed
			}
		} else {
			// else update dislike to like and change the count
			if existingReaction.Reaction_type == -1 {
				return nil
			} else {
				existingReaction.Reaction_type = -1
				if err := uc.blogReactionRepo.UpdateReaction(txCtx, existingReaction); err != nil {
					return domain.ErrUpdateBlogReactionFailed
				}
				// -1 like_count and +1 dislike_count
				if err := uc.blogRepo.ToggleLikeDislikeCounts(txCtx, blog_id, -1, 1); err != nil {
					return domain.ErrToggleLikeDislikeFailed
				}
			}
		}
		return nil
	})
}

func (uc *BlogReactionUseCase) UnlikeBlog(ctx context.Context, blogID, userID string) error {
	return uc.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		existingReaction, err := uc.blogReactionRepo.GetReactionByUserAndBlog(txCtx, blogID, userID)
		if err != nil {
			// desired state
			if errors.Is(err, domain.ErrBlogReactionNotFound) {
				return nil
			}
			return domain.ErrCheckBlogReactionFailed
		}

		// Delete the reaction
		if err := uc.blogReactionRepo.DeleteReaction(txCtx, blogID, userID); err != nil {
			return domain.ErrDeletingBlogReaction
		}

		// Only decrement if the existing reaction was a like
		if existingReaction.Reaction_type == 1 {
			if err := uc.blogRepo.DecrementLike(txCtx, blogID); err != nil {
				return domain.ErrToggleLikeDislikeFailed
			}
		}
		return nil
	})
}

func (uc *BlogReactionUseCase) UndislikeBlog(ctx context.Context, blogID, userID string) error {
	return uc.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		existingReaction, err := uc.blogReactionRepo.GetReactionByUserAndBlog(txCtx, blogID, userID)
		if err != nil {
			// No reaction exists — nothing to undo
			if errors.Is(err, domain.ErrBlogReactionNotFound) {
				return nil
			}
			return domain.ErrCheckBlogReactionFailed
		}

		// delete reaction
		if err := uc.blogReactionRepo.DeleteReaction(txCtx, blogID, userID); err != nil {
			return domain.ErrDeletingBlogReaction
		}

		// decrement if it was dislike
		if existingReaction.Reaction_type == -1 {
			if err := uc.blogRepo.DecrementDisLike(txCtx, blogID); err != nil {
				return domain.ErrToggleLikeDislikeFailed
			}
		}

		return nil
	})
}

package domain

import "errors"

var (
	// Blog-Specific Conditions
	ErrBlogNotFound        = errors.New("blog not found")
	ErrInvalidBlogID       = errors.New("invalid blog ID")
	ErrInvalidBlogIdFormat = errors.New("invalid blog ID format")
	ErrNoBlogChangesMade   = errors.New("no changes were made to the blog")

	// Repository-Specific Errors
	ErrQueryFailed         = errors.New("failed to execute MongoDB query")
	ErrDocumentDecoding    = errors.New("failed to decode MongoDB document")
	ErrCursorFailed        = errors.New("cursor encountered an error during iteration")
	ErrInsertingDocuments  = errors.New("failed to insert document(s)")
	ErrRetrievingDocuments = errors.New("failed to retrieve documents")
	ErrDecodingDocument    = errors.New("failed to decode document")
	ErrUpdatingDocument    = errors.New("failed to update document")
	ErrDeletingDocument    = errors.New("failed to delete document")
	ErrCursorIteration     = errors.New("database cursor iteration error")

	// Usecase-Specific Errors
	ErrBlogRequired   = errors.New("blog cannot be nil")
	ErrEmptyTitle     = errors.New("title is required")
	ErrEmptyContent   = errors.New("content is required")
	ErrInvalidUserID  = errors.New("user_id is required")
	ErrBlogIDRequired = errors.New("blog ID is required")

	// Blog-Reaction Related Errors
	ErrBlogReactionNotFound     = errors.New("blog reaction not found")
	ErrCheckBlogReactionFailed  = errors.New("failed to check existing blog reaction")
	ErrInternalServerError      = errors.New("internal server error")
	ErrCreateBlogReactionFailed = errors.New("failed to create blog reaction")
	ErrUpdateBlogReactionFailed = errors.New("failed to update blog reaction")
	ErrIncrementLikeFailed      = errors.New("failed to increment like count")
	ErrToggleLikeDislikeFailed  = errors.New("failed to toggle like/dislike counts")
	ErrDeletingBlogReaction     = errors.New("failed to delete blog reaction")
	ErrViewRecordAlreadyExists  = errors.New("view record already exists")
	ErrCreateViewRecordFailed   = errors.New("failed to create view record")
	ErrIncrementViewFailed      = errors.New("failed to increment blog view count")
)

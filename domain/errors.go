package domain

import "errors"

var (
	// Blog-Specific Conditions
	ErrBlogNotFound      = errors.New("blog not found")
	ErrInvalidBlogID     = errors.New("invalid blog ID")
	ErrNoBlogChangesMade = errors.New("no changes were made to the blog")

	// Comment-Specific Conditions
	ErrCommentNotFound      = errors.New("comment not found")
	ErrInvalidCommentID     = errors.New("invalid comment ID")
	ErrNoCommentChangesMade = errors.New("no changes were made to the comment")

	// Usecase-Specific Errors
	ErrBlogRequired      = errors.New("blog cannot be nil")
	ErrEmptyTitle        = errors.New("title is required")
	ErrEmptyContent      = errors.New("content is required")
	ErrInvalidUserID     = errors.New("user_id is required")
	ErrBlogIDRequired    = errors.New("blog ID is required")

	// Comment Usecase-Specific Errors
	ErrCommentRequired   = errors.New("comment cannot be nil")
	ErrEmptyCommentContent = errors.New("comment content is required")
	ErrCommentIDRequired = errors.New("comment ID is required")

	// Repository-Specific Errors
	ErrInsertingDocuments  = errors.New("failed to insert document(s)")
	ErrRetrievingDocuments = errors.New("failed to retrieve documents")
	ErrDecodingDocument    = errors.New("failed to decode document")
	ErrUpdatingDocument    = errors.New("failed to update document")
	ErrDeletingDocument    = errors.New("failed to delete document")
	ErrCursorIteration     = errors.New("database cursor iteration error")
)

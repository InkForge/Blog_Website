package domain

import "errors"

var (
	// Blog-Specific Conditions
	ErrBlogNotFound      = errors.New("blog not found")
	ErrInvalidBlogID     = errors.New("invalid blog ID")
	ErrNoBlogChangesMade = errors.New("no changes were made to the blog")

	// Repository-Specific Errors
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
)

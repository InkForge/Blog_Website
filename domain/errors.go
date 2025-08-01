package domain

import "errors"

var (
	// Blog-Specific Conditions
	ErrBlogNotFound        = errors.New("blog not found")
	ErrInvalidBlogID       = errors.New("invalid blog ID")
	ErrInvalidBlogIdFormat = errors.New("invalid blog ID format")
	ErrNoBlogChangesMade   = errors.New("no changes were made to the blog")

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

	ErrQueryFailed      = errors.New("failed to execute MongoDB query")
	ErrDocumentDecoding = errors.New("failed to decode MongoDB document")
	ErrCursorFailed     = errors.New("cursor encountered an error during iteration")

	ErrInsertingDocuments  = errors.New("failed to insert document(s)")
	ErrRetrievingDocuments = errors.New("failed to retrieve documents")
	ErrDecodingDocument    = errors.New("failed to decode document")
	ErrUpdatingDocument    = errors.New("failed to update document")
	ErrDeletingDocument    = errors.New("failed to delete document")
	ErrCursorIteration     = errors.New("database cursor iteration error")


	
	// User-Specific Errors
	ErrInvalidEmailFormat      = errors.New("invalid email format")
    ErrEmailAlreadyExists      = errors.New("email already exists")
    ErrUserNotFound            = errors.New("user not found")
    ErrInvalidCredentials      = errors.New("invalid credentials")
    ErrEmailNotVerified        = errors.New("email not verified")
    ErrPasswordHashingFailed   = errors.New("password hashing failed")
    ErrTokenGenerationFailed   = errors.New("token generation failed")
    ErrEmailSendingFailed      = errors.New("email sending failed")
    ErrUserCreationFailed      = errors.New("user creation failed")
    ErrDatabaseOperationFailed = errors.New("database operation failed")

)

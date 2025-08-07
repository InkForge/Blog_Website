package domain

import "errors"

var (
	// ─── Blog-Specific Errors ───────────────────────────────────────────────
	ErrBlogNotFound        = errors.New("blog not found")
	ErrInvalidBlogID       = errors.New("invalid blog ID")
	ErrInvalidBlogIdFormat = errors.New("invalid blog ID format")
	ErrNoBlogChangesMade   = errors.New("no changes were made to the blog")
	ErrBlogRequired        = errors.New("blog cannot be nil")
	ErrEmptyTitle          = errors.New("title is required")
	ErrEmptyContent        = errors.New("content is required")
	ErrInvalidUserID       = errors.New("user_id is required")
	ErrBlogIDRequired      = errors.New("blog ID is required")
	ErrIncrementViewFailed = errors.New("failed to increment blog view count")
	ErrNotBlogAuthor       = errors.New("user is not the author of the blog")

	// ─── Blog Reaction Errors ──────────────────────────────────────────────
	ErrBlogReactionNotFound     = errors.New("blog reaction not found")
	ErrCheckBlogReactionFailed  = errors.New("failed to check existing blog reaction")
	ErrCreateBlogReactionFailed = errors.New("failed to create blog reaction")
	ErrUpdateBlogReactionFailed = errors.New("failed to update blog reaction")
	ErrIncrementLikeFailed      = errors.New("failed to increment like count")
	ErrToggleLikeDislikeFailed  = errors.New("failed to toggle like/dislike counts")
	ErrDeletingBlogReaction     = errors.New("failed to delete blog reaction")

	// ─── Blog View Errors ──────────────────────────────────────────────────
	ErrViewRecordAlreadyExists = errors.New("view record already exists")
	ErrCreateViewRecordFailed  = errors.New("failed to create view record")

	// ─── Comment Errors ───────────────────────────────────────────────────
	ErrCommentNotFound      = errors.New("comment not found")
	ErrInvalidCommentID     = errors.New("invalid comment ID")
	ErrNoCommentChangesMade = errors.New("no changes were made to the comment")
	ErrCommentRequired      = errors.New("comment cannot be nil")
	ErrEmptyCommentContent  = errors.New("comment content is required")
	ErrCommentIDRequired    = errors.New("comment ID is required")
	ErrForbidden            = errors.New("forbidden")

	// ─── Comment Reaction Errors ──────────────────────────────────────────
	ErrCommentReactionNotFound = errors.New("comment reaction not found")
	ErrInvalidReactionAction   = errors.New("invalid reaction action")
	ErrReactionAlreadyExists   = errors.New("reaction already exists for this user and comment")
	ErrCommentReactionRequired = errors.New("comment reaction cannot be nil")
	ErrInvalidReactionType     = errors.New("invalid reaction type")

	// ─── Repository-Level Errors ──────────────────────────────────────────
	ErrQueryFailed         = errors.New("failed to execute MongoDB query")
	ErrDocumentDecoding    = errors.New("failed to decode MongoDB document")
	ErrCursorFailed        = errors.New("cursor encountered an error during iteration")
	ErrInsertingDocuments  = errors.New("failed to insert document(s)")
	ErrRetrievingDocuments = errors.New("failed to retrieve documents")
	ErrDecodingDocument    = errors.New("failed to decode document")
	ErrUpdatingDocument    = errors.New("failed to update document")
	ErrDeletingDocument    = errors.New("failed to delete document")
	ErrCursorIteration     = errors.New("database cursor iteration error")

	// ─── User Errors ───────────────────────────────────────────────────────
	ErrInvalidToken                     = errors.New("invalid token")
	ErrTokenRevocationFailed            = errors.New("token revocation failed")
	ErrInvalidEmailFormat               = errors.New("invalid email format")
	ErrEmailAlreadyExists               = errors.New("email already exists")
	ErrUserNotFound                     = errors.New("user not found")
	ErrInvalidCredentials               = errors.New("invalid credentials")
	ErrEmailNotVerified                 = errors.New("email not verified")
	ErrPasswordHashingFailed            = errors.New("password hashing failed")
	ErrTokenGenerationFailed            = errors.New("token generation failed")
	ErrEmailSendingFailed               = errors.New("email sending failed")
	ErrUserCreationFailed               = errors.New("user creation failed")
	ErrDatabaseOperationFailed          = errors.New("database operation failed")
	ErrOAuthUserCannotLoginWithPassword = errors.New("OAuth user cannot login with password")
	ErrInvalidInput                     = errors.New("invalid input")
	ErrUserUpdateFailed                 = errors.New("user update failed")
	ErrEmailVerficationFailed           = errors.New("email verification failed")
	ErrTokenVerificationFailed          = errors.New("token verification failed")
	ErrInvalidRole                      = errors.New("invalid role")
	ErrWeakPassword                     = errors.New("weak password")
	ErrPasswordMismatch                 = errors.New("passwords do not match")
	ErrUserVerified                     = errors.New("user already verified")
	ErrGetTokenExpiryFailed             = errors.New("failed to get token expiration time")
	ErrWeakPassword                     = errors.New("password is too weak")
	ErrInvalidRole                      = errors.New("invalid role specified")
	ErrInvalidOAuthUserData             = errors.New("invalid OAuth user data")
	ErrOAuthProviderMismatch             = errors.New("OAuth provider mismatch for this account")

	// ─── Generic Errors ────────────────────────────────────────────────────
	ErrInternalServerError = errors.New("internal server error")
)

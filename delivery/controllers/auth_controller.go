package controllers

import (
	"errors"
	"net/http"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthUsecase domain.IAuthUsecase
}

// ChangePassword handles password change requests.
// It requires the user to be authenticated and to provide the old and new passwords.
// Returns 400 for invalid input, 401 if user is unauthenticated or old password is wrong.
func (au *AuthController) ChangePassword(c *gin.Context) {
	type payload struct {
		OldPassword string `json:"old_password" binding:"required;min=6"`
		NewPassword string `json:"new_password" binding:"required;min=6"`
	}

	var body payload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to fetch userID"})
		return
	}

	err := au.AuthUsecase.ChangePassword(c, userID, body.OldPassword, body.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		case errors.Is(err, domain.ErrPasswordMismatch):
			c.JSON(http.StatusBadRequest, gin.H{"error": "old password is not correct"})
		case errors.Is(err, domain.ErrInvalidCredentials):
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update to new password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": " failed to update password", "details": err.Error()})
		}
	}

}

// ResetPassword handles password reset using a reset token sent via email.
// It expects a valid token and a new password, and updates the user’s password if valid.
// Returns 400 for invalid token or password input, or 500 on internal failure.
func (au *AuthController) ResetPassword(c *gin.Context) {
	type payload struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required;min=6"`
	}

	var body payload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	err := au.AuthUsecase.ResetPassword(c, body.Token, body.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		case errors.Is(err, domain.ErrInvalidToken):
			c.JSON(http.StatusBadRequest, gin.H{"error": "token is malformed or expired", "details": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update the user", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})

}
func (au *AuthController) RequestPasswordReset(c *gin.Context) {

	type payload struct {
		Email string `json:"email" binding:"required"`
	}

	var body payload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	err := au.AuthUsecase.RequestPasswordReset(c, body.Email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		case errors.Is(err, domain.ErrEmailSendingFailed):
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to send email"})
		case errors.Is(err, domain.ErrInvalidToken):
			c.JSON(http.StatusBadRequest, gin.H{"error": "token is malformed or expired", "details": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send verification email", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset link sent successfully"})

}

// ResendVerification handles resending a verification email to a user who is not yet verified.
// Accepts an email, verifies its format and status, then sends a new verification link.
func (au *AuthController) ResendVerification(c *gin.Context) {
	type payload struct {
		Email string `json:"email" binding:"required"`
	}

	var body payload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	err := au.AuthUsecase.ResendVerificationEmail(c, body.Email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		case errors.Is(err, domain.ErrEmailSendingFailed):
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to send email"})
		case errors.Is(err, domain.ErrUserVerified):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already verified"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resend verification email", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification link sent successfully"})

}

// VerifyEmail handles verification of a user’s email via a token sent in the URL.
// It validates the token and updates the user's verified status if valid.
func (au *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token missing"})
		return
	}

	err := au.AuthUsecase.VerifyEmail(c, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to verify email", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user verified successfully"})
}

// RefreshToken is an HTTP handler that handles the token refreshing endpoint.
// If the operation is successful, it updates the access_token cookie and returns its expiration.
func (au *AuthController) RefreshToken(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to fetch userID"})
		return
	}

	accessToken, _, expiresIn, err := au.AuthUsecase.RefreshToken(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    *accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(expiresIn.Seconds()),
	})

	resp := map[string]interface{}{
		"expires_in": int(expiresIn.Seconds()),
	}
	c.JSON(http.StatusOK, resp)

}

// Logout is an HTTP handler that handles user logout and clears session-related cookies.
// It also performs any backend cleanup such as session invalidation.
func (au *AuthController) Logout(c *gin.Context) {
	// get the user id
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to fetch userID"})
		return
	}

	err := au.AuthUsecase.Logout(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log user out", "details": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

package controllers

import (
	"net/http"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

type OAuth2Controller struct {
	OAuthService domain.IOAuth2Service
	AuthUsecase  domain.IAuthUsecase
}

func NewOAuth2Controller(service domain.IOAuth2Service, authUsecase domain.IAuthUsecase) *OAuth2Controller {
	return &OAuth2Controller{
		OAuthService: service,
		AuthUsecase:  authUsecase,
	}
}

func (ctrl *OAuth2Controller) RedirectToProvider(c *gin.Context) {
	provider := c.Param("provider")
	state := "random-state" 
	url, err := ctrl.OAuthService.GetAuthorizationURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, url)
}

func (ctrl *OAuth2Controller) HandleCallback(c *gin.Context) {
	ctx := c.Request.Context()
	provider := c.Param("provider")
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code in query"})
		return
	}

	// authenticate with the provider
	oauthUser, err := ctrl.OAuthService.Authenticate(ctx, provider, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// register/login via usecase
	result, err := ctrl.AuthUsecase.OAuthLogin(ctx, oauthUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// set auth token cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    result.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(result.ExpiresIn.Seconds()),
	})

	// return safe user response
	safeUser := gin.H{
		"user_id":   result.User.UserID,
		"email":     result.User.Email,
		"username":  result.User.Username,
		"firstName": result.User.FirstName,
		"lastName":  result.User.LastName,
		"role":      result.User.Role,
		"provider":  result.User.Provider,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OAuth login successful",
		"user":    safeUser,
	})
}

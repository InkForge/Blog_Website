package infrastructures

import (
	"fmt"
	"net/http"

	"github.com/InkForge/Blog_Website/domain"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	jwtSecret []byte
	jwtService domain.IJWTService
}

func NewAuthService(jwtService domain.IJWTService, secret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(secret),
		jwtService: jwtService,
	}
}

// define constants for claim keys
const (
	ClaimUserID   = "sub"
	ClaimUserRole = "role"
)

func (a *AuthService) AuthWithRole(allowedRoles ...string) gin.HandlerFunc{
	return func(c *gin.Context){
		//read the cookie
		cookie,err:=c.Request.Cookie("auth_token")
		if err!=nil{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"unauthorized: no auth cookie"})
			return 
		}
		tokenStr:=cookie.Value

		userID, userRole, err := a.jwtService.ValidateAccessToken(tokenStr)
		if err!=nil{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":fmt.Sprintf("unauthorized:invalid token(%v)",err)})
			return 
		}

		//save user data to context
		c.Set("userID",userID)
		c.Set("userRole",userRole)

		//check if role is authorzied
		authorized:=false
		for _,r := range allowedRoles{
			if userRole==r{
				authorized=true
				break
			}
		}
		if !authorized{
			c.AbortWithStatusJSON(http.StatusForbidden,gin.H{"error":"forbidden: role not authorized"})
			return
		}
		c.Next()

	}}
	

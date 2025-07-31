package infrastructures

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(secret)}
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

		//parse the token
		token,err:=jwt.Parse(tokenStr,func(token *jwt.Token)(interface {},error){
			if _,ok:=token.Method.(*jwt.SigningMethodHMAC);!ok{
				return nil,fmt.Errorf("unexpected signing method")
			}
			return a.jwtSecret,nil
		})
		//handle parse errors or invalid token
		if err!=nil || !token.Valid{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":fmt.Sprintf("unauthorized:invalid token(%v)",err)})
			return 
		}

		//extract claims
		claims,ok:=token.Claims.(jwt.MapClaims)
		if !ok{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"unauthorized: invalid token claims"})
			return 
		}

		userID,ok1:=claims[ClaimUserID].(string)
		role,ok2:=claims[ClaimUserRole].(string)

		if !ok1 || !ok2{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"unauthorized:missing user infor in token"})
			return 
		}

		//save user data to context
		c.Set("userID",userID)
		c.Set("userRole",role)

		//check if role is authorzied
		authorized:=false
		for _,r := range allowedRoles{
			if role==r{
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
	

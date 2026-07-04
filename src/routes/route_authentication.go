package routes

import (
	"github.com/gin-gonic/gin"
	"gourmand.golang.arango/src/interfaces"
)

type AuthenticationRoute struct {
	AuthenticationService interfaces.IAuthentication
}

func NewAuthenticationRoute(authService interfaces.IAuthentication) *AuthenticationRoute {
	return &AuthenticationRoute{
		AuthenticationService: authService,
	}
}


func (r *AuthenticationRoute) GenerateToken(c *gin.Context, email, username, password string) (string, error) {
	ctx := c.Request.Context()
	return r.AuthenticationService.GenerateToken(ctx, email, password)
}

func (r *AuthenticationRoute) AuthenticateUser(c *gin.Context) (bool, error) {
	ctx := c.Request.Context()
	return r.AuthenticationService.AuthenticateUser(ctx, c.PostForm("email"), c.PostForm("password"))
}
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

func (r *AuthenticationRoute) Refresh(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.PostForm("email")

	token, err := r.AuthenticationService.GenerateToken(ctx, email)
	
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (r *AuthenticationRoute) AuthenticateUser(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := r.AuthenticationService.AuthenticateUser(ctx, c.PostForm("email"), c.PostForm("password"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !result {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := r.AuthenticationService.GenerateToken(ctx, c.PostForm("email"))
	
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (r *AuthenticationRoute) AuthenticateToken(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.GetHeader("Token")

	valid, err := r.AuthenticationService.AuthenticateToken(ctx, token)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !valid {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(200, gin.H{"valid": true})
}
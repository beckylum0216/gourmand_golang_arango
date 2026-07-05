package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/interfaces"
)

type UserRoutes struct {
	userService interfaces.IUser
}

func NewUserRoutes(userService interfaces.IUser) *UserRoutes {
	return &UserRoutes{
		userService: userService,
	}
}

func (ur *UserRoutes) CreateUser(c *gin.Context) {
	var user entities.User
	var person entities.Person

	var req struct {
		User   entities.User   `json:"user"`
		Person entities.Person `json:"person"`
		Auth   entities.Authentication `json:"auth"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user = req.User
	person = req.Person
	auth := req.Auth

	if err := ur.userService.CreateUser(c.Request.Context(), &person, &user, &auth); err != nil {
		// return a sensible error instead of letting a panic bubble up
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "User created successfully",
		"user":    user,
	})

}

func (ur *UserRoutes) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := ur.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}

func (ur *UserRoutes) GetUsers(c *gin.Context) {
	users, err := ur.userService.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(404, gin.H{"error": "No users found"})
		return
	}
	c.JSON(200, users)
}

func (ur *UserRoutes) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := ur.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

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
	if err := c.ShouldBindJSON(&user); 
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ur.userService.CreateUser(c.Request.Context(), &person, &user); 
	err != nil {
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

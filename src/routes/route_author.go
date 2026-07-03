package routes

import (
	"github.com/gin-gonic/gin"
	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/interfaces"
)

type AuthorRoutes struct {
	authorService interfaces.IAuthor
}

func NewAuthorRoutes(authorService interfaces.IAuthor) *AuthorRoutes {
	return &AuthorRoutes{authorService: authorService}
}

func (ar *AuthorRoutes) CreateAuthor(c *gin.Context) {
	// 1. Define a DTO that matches the incoming JSON structure
	var req struct {
		Author entities.Author `json:"author"`
		Person entities.Person `json:"person"`
	}
	
	ctx := c.Request.Context()
	// 2. Bind the JSON once
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON payload"})
		return
	}
	// 3. Pass the unmarshaled data to your service
	if err := ar.authorService.CreateAuthor(ctx, &req.Person, &req.Author); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Author created successfully",
		"author":  req.Author,
	})
}

func (ar *AuthorRoutes) GetAuthor(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	author, err := ar.authorService.GetAuthor(ctx, id)

	if err != nil {
		c.JSON(404, gin.H{"error": "Author not found"})
		return
	}
	c.JSON(200, author)
}

func (ar *AuthorRoutes) UpdateAuthor(c *gin.Context) {
	ctx := c.Request.Context()
	var authorData = &entities.Author{}
	id := c.Param("id")

	if err := c.ShouldBindJSON(authorData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := ar.authorService.UpdateAuthor(ctx, id, authorData); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Author updated successfully",
	})
}

func (ar *AuthorRoutes) DeleteAuthor(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	err := ar.authorService.DeleteAuthor(ctx, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Author not found"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Author deleted successfully",
	})
}

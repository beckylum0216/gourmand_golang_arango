package routes

import (

	"github.com/gin-gonic/gin"
	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/interfaces"
)

type RecipeRoutes struct {
	recipeService interfaces.IRecipe
}

func NewRecipeRoutes(recipeService interfaces.IRecipe) *RecipeRoutes {
	return &RecipeRoutes{recipeService: recipeService}
}

func (rr *RecipeRoutes) CreateRecipe(c *gin.Context) {
	var recipe entities.Recipe
	var author entities.Author
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := rr.recipeService.CreateRecipe(c.Request.Context(), &author, &recipe); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Recipe created successfully",
		"recipe":  recipe,
	})
}

func (rr *RecipeRoutes) TranscribeRecipe(c *gin.Context) {
	var recipe entities.Recipe
	var author entities.Author

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := rr.recipeService.TranscribeRecipe(c.Request.Context(), &author, &recipe); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Recipe created successfully",
		"recipe":  recipe,
	})
}

func (rr *RecipeRoutes) GetRecipe(c *gin.Context) {
	id := c.Param("id")

	recipe, err := rr.recipeService.GetRecipe(c.Request.Context(), id)

	if err != nil {
		c.JSON(404, gin.H{"error": "Recipe not found"})
		return
	}
	c.JSON(200, recipe)
}

func (rr *RecipeRoutes) GetRecipes(c *gin.Context) {
	recipes, err := rr.recipeService.GetRecipes(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve recipes"})
		return
	}
	c.JSON(200, recipes)
}

func (rr *RecipeRoutes) UpdateRecipe(c *gin.Context) {
	var recipeData = &entities.Recipe{}

	if err := c.ShouldBindJSON(&recipeData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	id := c.Param("id")
	if err := rr.recipeService.UpdateRecipe(c.Request.Context(), id, recipeData); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := rr.recipeService.UpdateRecipe(c.Request.Context(), id, recipeData); err != nil {
		c.JSON(200, gin.H{
			"message": "Recipe updated successfully",
			"recipe":  recipeData,
		})
	}
}

func (rr *RecipeRoutes) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	if err := rr.recipeService.DeleteRecipe(c.Request.Context(), id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Recipe deleted successfully"})
}
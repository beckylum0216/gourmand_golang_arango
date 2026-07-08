package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gourmand.golang.arango/src/database"
	"gourmand.golang.arango/src/routes"
	"gourmand.golang.arango/src/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, falling back to system environment variables")
	}

	log.Println("ARANGO_ENDPOINT:", os.Getenv("ARANGO_ENDPOINT"))
	log.Println("ARANGO_USER:", os.Getenv("ARANGO_USER"))
	log.Println("ARANGO_PASSWORD:", os.Getenv("ARANGO_PASSWORD"))
	log.Println("ARANGO_DB_NAME:", os.Getenv("ARANGO_DB_NAME"))

	ctx := context.Background()
	dsn := os.Getenv("ARANGO_ENDPOINT")
	userName := os.Getenv("ARANGO_USER")
	password := os.Getenv("ARANGO_PASSWORD")
	name := os.Getenv("ARANGO_DB_NAME")
	db, err := database.NewArangoDB(ctx, dsn, userName, password, name)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:4200",
			"http://127.0.0.1:4200",
			"http://[::1]:4200",
			"http://15.134.229.139:4200",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"http://localhost:4200", "http://127.0.0.1:4200", "http://[::1]:4200"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, MaxAge: 12 * time.Hour,
	}))

	userService := services.NewUserService(*db)
	userRoutes := routes.NewUserRoutes(userService)
	personService := services.NewPersonService(*db)
	personRoutes := routes.NewPersonRoutes(personService)
	authorService := services.NewAuthorService(*db)
	authorRoutes := routes.NewAuthorRoutes(authorService)
	recipeService := services.NewRecipeService(*db)
	recipeRoutes := routes.NewRecipeRoutes(recipeService)
	authService, err := services.NewAuthenticationService(*db)
	if err != nil {
		log.Fatalf("failed to initialize auth service: %v", err)
	}
	authRoutes := routes.NewAuthenticationRoute(authService)

	api := router.Group("/api")
	{
		api.POST("/create_user", userRoutes.CreateUser)
		api.GET("/get_user/:id", userRoutes.GetUser)
		api.GET("/get_users", userRoutes.GetUsers)
		api.POST("/delete_user/:id", userRoutes.DeleteUser)

		api.GET("/get_person/:id", personRoutes.GetPerson)
		api.GET("/get_persons", personRoutes.GetPersons)
		api.POST("/delete_person/:id", personRoutes.DeletePerson)

		api.POST("/create_author", authorRoutes.CreateAuthor)
		api.GET("/get_author/:id", authorRoutes.GetAuthor)
		api.GET("/get_authors", authorRoutes.GetAuthors)
		api.POST("/delete_author/:id", authorRoutes.DeleteAuthor)

		api.POST("/create_recipe", recipeRoutes.CreateRecipe)
		api.POST("/transcribe_recipe", recipeRoutes.TranscribeRecipe)
		api.GET("/get_recipe/:id", recipeRoutes.GetRecipe)
		api.GET("/get_recipes", recipeRoutes.GetRecipes)
		api.POST("/delete_recipe/:id", recipeRoutes.DeleteRecipe)

		api.POST("/login", authRoutes.AuthenticateUser)
		api.GET("/authenticate_token", authRoutes.AuthenticateToken)
	}

	router.Run("0.0.0.0:8080")
}

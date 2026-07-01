package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"errors"

	"github.com/ollama/ollama/api"
	"github.com/arangodb/go-driver/v2/arangodb"

	"gourmand.golang.arango/src/entities"
	"gourmand.golang.arango/src/interfaces"
)

const recipes_collection = "recipes"
const authors_recipes_edges = "authors_recipes"

type RecipeService struct {
	db arangodb.Database
}

func NewRecipeService(db arangodb.Database) interfaces.IRecipe {
	return &RecipeService{db: db}
}

func (s *RecipeService) CreateRecipe(ctx context.Context, author *entities.Author, recipe *entities.Recipe) error {
	col, err := s.db.GetCollection(ctx, recipes_collection, nil)
	if err != nil {
		return err
	}

	meta, err := col.CreateDocument(ctx, recipe)
	if err != nil {
		return err
	}

	recipe.Id = meta.Key

	edgeCol, err := s.db.GetCollection(ctx, authors_recipes_edges, nil)
	if err != nil {
		return err
	}

	edge := map[string]interface{}{
		"_from": "authors/" + author.Id,
		"_to":   "recipes/" + recipe.Id,
	}

	_, err = edgeCol.CreateDocument(ctx, edge)
	if err != nil {
		return err
	}

	if _, err := edgeCol.CreateDocument(ctx, edge); 
	err != nil {
		return errors.New("failed to create edge between author and recipe: " + err.Error())
	}

	return nil
}

func (s *RecipeService) TranscribeRecipe(ctx context.Context, 
	author *entities.Author, recipe *entities.Recipe) error {
	ollamaURL, _ := url.Parse("http://192.168.0.222:11434")
	client := api.NewClient(ollamaURL, http.DefaultClient)

	prompt := `You are a recipe transcription tool. 
	Convert ONLY the information provided below into markdown. 
	Do NOT add ingredients or steps that aren't mentioned.

		Format:
		---
		title: Recipe Name
		yield: X servings
		prep_time: Xmin
		cook_time: Xmin
		tags: [tag1, tag2]
		---

		## Ingredients
		- ingredient list

		## Instructions
		1. numbered steps

		Recipe text:
		` + recipe.Method + ``

	req := &api.GenerateRequest{
		Model:  "llama3.2:3b",
		Prompt: prompt,
		Options: map[string]interface{}{
			"temperature": 0.1,
		},
	}

	var markdown strings.Builder

	response := func(resp api.GenerateResponse) error {
		markdown.WriteString(resp.Response)
		return nil
	}

	err := client.Generate(ctx, req, response)
	if err != nil {
		return fmt.Errorf("ollama generate failed: %w", err)
	}

	recipe.RecipeText = markdown.String()

	col, err := s.db.GetCollection(ctx, recipes_collection, nil)
	if err != nil {
		return err
	}

	meta, err := col.CreateDocument(ctx, recipe)
	if err != nil {
		return err
	}

	recipe.Id = meta.Key

	edgeCol, err := s.db.GetCollection(ctx, authors_recipes_edges, nil)
	if err != nil {
		return err
	}

	edge := map[string]interface{}{
		"_from": "authors/" + author.Id,
		"_to":   "recipes/" + recipe.Id,
	}

	if _, err := edgeCol.CreateDocument(ctx, edge); 
	err != nil {
		return errors.New("failed to create edge between author and recipe: " + err.Error())
	}

	return nil
}

func (s *RecipeService) GetRecipe(ctx context.Context, id string) (*entities.Recipe, error) {
	var recipe entities.Recipe
	col, err := s.db.GetCollection(ctx, "recipes", nil)
	if err != nil {
		return nil, err
	}

	meta, err := col.ReadDocument(ctx, id, &recipe)
	if err != nil {
		return nil, err
	}

	recipe.Id = meta.Key

	return &recipe, nil
}

func (s *RecipeService) GetRecipes(ctx context.Context) ([]*entities.Recipe, error) {
	var recipes []*entities.Recipe
	
	query := "FOR r IN @@recipes_collection RETURN r"
	bindVars := map[string]interface{}{
		"@recipes": recipes_collection,
	}

	options := &arangodb.QueryOptions{
		BindVars: bindVars,
	}

	cursor, err := s.db.Query(ctx, query, options)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for cursor.HasMore() {
		var recipe entities.Recipe
		_, err := cursor.ReadDocument(ctx, &recipe)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}



	return recipes, nil
}

func (s *RecipeService) UpdateRecipe(ctx context.Context, id string, recipe *entities.Recipe) error {
	col, err := s.db.GetCollection(ctx, recipes_collection, nil)
	if err != nil {
		return err
	}

	patch := map[string]interface{}{
		"title":        recipe.Title,
		"description":  recipe.Description,
		"prep_time":    recipe.PrepTime,
		"recipe_text":  recipe.RecipeText,
	}

	_, err = col.UpdateDocument(ctx, id, patch)
	if err != nil {
		return err
	}

	return nil
}

func (s *RecipeService) DeleteRecipe(ctx context.Context, id string) error {
	col, err := s.db.GetCollection(ctx, recipes_collection, nil)
	if err != nil {
		return err
	}

	options := arangodb.CollectionDocumentDeleteOptions{
		IfMatch: "",
	}

	_, err = col.DeleteDocumentWithOptions(ctx, id, &options)
	if err != nil {
		return err
	}


	return err
}

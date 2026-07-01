package interfaces

import (
	"context"

	"gourmand.golang.arango/src/entities"
)

type IRecipe interface {
	CreateRecipe(ctx context.Context, author *entities.Author, recipe *entities.Recipe) error
	TranscribeRecipe(ctx context.Context, author *entities.Author, recipe *entities.Recipe) error
	GetRecipe(ctx context.Context, id string) (*entities.Recipe, error)
	GetRecipes(ctx context.Context) ([]*entities.Recipe, error)
	UpdateRecipe(ctx context.Context, id string, recipe *entities.Recipe) error
	DeleteRecipe(ctx context.Context, id string) error
}

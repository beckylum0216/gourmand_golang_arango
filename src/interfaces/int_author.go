package interfaces

import (
	"context"

	"gourmand.golang.arango/src/entities"
)

type IAuthor interface {
	CreateAuthor(ctx context.Context, person *entities.Person, author *entities.Author) error
	GetAuthor(ctx context.Context, id string) (*entities.Author, error)
	GetAuthors(ctx context.Context) ([]*entities.Author, error)
	UpdateAuthor(ctx context.Context, id string, author *entities.Author) error
	DeleteAuthor(ctx context.Context, id string) error
}

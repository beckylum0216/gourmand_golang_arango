package interfaces

import (
	"context"
	"gourmand.golang.arango/src/entities"
)

type IPerson interface {
	CreatePerson(ctx context.Context, person *entities.Person) (*entities.Person, error)
	GetPerson(ctx context.Context, id string) (*entities.Person, error)
	GetPersons(ctx context.Context) ([]*entities.Person, error)
	UpdatePerson(ctx context.Context, id string, person *entities.Person) error
	DeletePerson(ctx context.Context, id string) error
}

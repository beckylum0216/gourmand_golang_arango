package interfaces

import (
	"context"
	"gourmand.golang.arango/src/entities"
)

type IUser interface {
	CreateUser(ctx context.Context, person *entities.Person, user *entities.User, auth *entities.Authentication) error
	GetUser(ctx context.Context, id string) (*entities.User, error)
	GetUsers(ctx context.Context) ([]*entities.User, error)
	UpdateUser(ctx context.Context, id string, user *entities.User) error
	DeleteUser(ctx context.Context, id string) error
}
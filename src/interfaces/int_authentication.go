package interfaces

import (
	"context"

	"gourmand.golang.arango/src/entities"
)

type IAuthentication interface {
	CreateAuthentication(ctx context.Context, email, password string) (*entities.Authentication, error)
	AuthenticateUser(ctx context.Context, email, password string) (bool, error)
	GenerateToken(ctx context.Context, email string) (string, error)
	AuthenticateToken(ctx context.Context, token string) (bool, error)
}

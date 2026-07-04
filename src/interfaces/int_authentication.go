package interfaces

import (
	"context"
)

type IAuthentication interface {
	CreateAuthentication(ctx context.Context, email, password string) error
	AuthenticateUser(ctx context.Context, email, password string) (bool, error)
	GenerateToken(ctx context.Context, email, password string) (string, error)
}
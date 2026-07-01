package interfaces

import (
	"context"
)

type IAuthentication interface {
	AuthenticateUser(ctx context.Context, email, password string) (bool, error)
	GenerateToken(ctx context.Context, email, username, password string) (string, error)
}
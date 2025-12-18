package auth

import (
	"context"
	"main/internal/entity"
)

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

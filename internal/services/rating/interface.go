package rating

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data Create, customerId int64) (int64, error)
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

package cart

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, productId, customerId int64) (int64, error)
	Update(ctx context.Context, cartItemId int64, customerId int64) error
	DeleteCartItem(ctx context.Context, cartItemId int64, customerId int64) error
	GetList(ctx context.Context, filter entity.Filter) ([]Get, int64, error)
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

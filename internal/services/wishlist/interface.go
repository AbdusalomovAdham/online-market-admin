package wishlist

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	GetList(ctx context.Context, userId int64) ([]GetList, int64, error)
	Create(ctx context.Context, productId Create, userId int64) (int64, error)
	Delete(ctx context.Context, productId, userId int64) error
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

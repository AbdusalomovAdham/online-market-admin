package order

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, order Create, userId int64) error
	GetList(ctx context.Context, userId int64) ([]Get, error)
	GetById(ctx context.Context, orderId, userId int64) (Get, error)
	Delete(ctx context.Context, orderId int64, userId int64) error
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

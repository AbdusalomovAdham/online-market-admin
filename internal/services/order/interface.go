package order

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, order Create, userId int64) error
	GetList(ctx context.Context, userId int64, filter entity.Filter) ([]GetList, int64, error)
	GetById(ctx context.Context, orderId int64, filter entity.Filter) (Get, error)
	Delete(ctx context.Context, orderId int64, userId int64) error
	Update(ctx context.Context, updateData Update, orderId int64, adminId int64) error
	DeleteOrderItem(ctx context.Context, orderId, adminId int64) error
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

package order_status

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, orderStatus Create, userId int64) (int64, error)
	Delete(ctx context.Context, id int64, userId int64) error
	GetById(ctx context.Context, id int64) (OrderStatusById, error)
	GetList(ctx context.Context, filter entity.Filter, lang string) ([]Get, int, error)
	Update(ctx context.Context, id int64, data Update, userId int64) error
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

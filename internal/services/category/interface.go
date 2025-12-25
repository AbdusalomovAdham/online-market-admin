package category

import (
	"context"
	"main/internal/entity"
)

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

type Repository interface {
	Create(ctx context.Context, category Create, userId int64) (int64, error)
	Delete(ctx context.Context, id int64, userId int64) error
	GetById(ctx context.Context, id int64) (CategoryById, error)
	GetList(ctx context.Context, filter entity.Filter) ([]Get, int, error)
	Update(ctx context.Context, id int64, data Update, userId int64) error
}

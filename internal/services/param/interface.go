package param

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, param Create, adminId int64) (int64, error)
	GetById(ctx context.Context, paramId int64) (ParamById, error)
	GetList(ctx context.Context, filter entity.Filter) ([]Get, int, error)
	Delete(ctx context.Context, paramId int64, adminId int64) error
	Update(ctx context.Context, paramId int64, data UpdateParam, adminId int64) error
	GetByCategoryId(ctx context.Context, categoryId int64, filter entity.Filter) ([]GetByCategoryId, int64, error)
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

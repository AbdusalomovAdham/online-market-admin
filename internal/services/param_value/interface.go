package param_value

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, param Create, adminId int64) (int64, error)
	GetList(ctx context.Context, filter entity.Filter) ([]Get, int, error)
	Delete(ctx context.Context, paramValueId int64, adminId int64) error
	GetById(ctx context.Context, paramId int64) (ParamValueById, error)
	Update(ctx context.Context, id int64, data Update, userId int64) error
	GetListByParamId(ctx context.Context, filter entity.Filter, paramId int) ([]Get, error)
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

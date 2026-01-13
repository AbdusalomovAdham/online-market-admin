package district

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Get(ctx context.Context, filter entity.Filter) ([]GetList, int, error)
	GetListByRegionId(ctx context.Context, filter entity.Filter, regionId int) ([]GetList, error)
}

type Auth interface{}

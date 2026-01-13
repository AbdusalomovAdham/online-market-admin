package region

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Get(ctx context.Context, filter entity.Filter) ([]GetList, int, error)
}

type Auth interface{}

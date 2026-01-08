package role

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	GetList(ctx context.Context, filter entity.Filter) ([]Get, int, error)
}

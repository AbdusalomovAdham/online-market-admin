package auth

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	GetById(ctx context.Context, id int) (entity.User, error)
}

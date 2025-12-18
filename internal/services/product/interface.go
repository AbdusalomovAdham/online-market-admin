package products

import (
	"context"
	"main/internal/entity"
	"mime/multipart"
)

type Repository interface {
	Create(ctx context.Context, data Create, userId int64) (int64, error)
	GetById(ctx context.Context, id int64) (Get, error)
	GetList(ctx context.Context, filter entity.Filter, userId *int64) ([]Get, int, error)
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

type File interface {
	MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]entity.File, error)
}

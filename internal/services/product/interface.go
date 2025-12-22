package product

import (
	"context"
	"main/internal/entity"
	"mime/multipart"
)

type Repository interface {
	Create(ctx context.Context, data Create, userId int64) (int64, error)
	GetById(ctx context.Context, id int64) (GetById, error)
	GetByIdDetail(ctx context.Context, id int64) (GetById, error)
	GetList(ctx context.Context, filter entity.Filter) ([]Get, int, error)
	UpdateProduct(ctx context.Context, productId int, data Update, userId int64) error
	Delete(ctx context.Context, productId, userId int64) error
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

type File interface {
	MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string, startID *int32) ([]entity.File, error)
	Delete(ctx context.Context, url string) error
}

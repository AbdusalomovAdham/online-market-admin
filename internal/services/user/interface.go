package user

import (
	"context"
	"main/internal/entity"
	"main/internal/usecase/auth"
	"mime/multipart"
	"time"
)

type Repository interface {
	Create(ctx context.Context, data Create, adminId int64, birthTime time.Time) (int64, error)
	GetAll(ctx context.Context, filter entity.Filter) ([]Get, int64, error)
	GetById(ctx context.Context, id int64) (Get, error)
	Update(ctx context.Context, id int64, data Update, adminId int64) error
	Delete(ctx context.Context, id, adminId int64) error
}

type Auth interface {
	GenerateToken(ctx context.Context, data auth.GenerateToken) (string, error)
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateResetToken(n int) (string, error)
}

type File interface {
	MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string, startID *int32) ([]entity.File, error)
	Upload(ctx context.Context, image *multipart.FileHeader, folder string) (entity.File, error)
	Delete(ctx context.Context, url string) error
}

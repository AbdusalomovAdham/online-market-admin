package products

import (
	"context"
	"main/internal/entity"
	"mime/multipart"
)

type Service struct {
	repo Repository
	auth Auth
	file File
}

func NewService(repo Repository, auth Auth, file File) Service {
	return Service{repo: repo, auth: auth, file: file}
}

func (s Service) CreateProduct(ctx context.Context, data Create, authHeader string) (int64, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, data, isValidToken.Id)
}

func (s Service) MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]entity.File, error) {
	return s.file.MultipleUpload(ctx, files, folder)
}

func (s Service) GetById(ctx context.Context, id int64) (Get, error) {
	data, err := s.repo.GetById(ctx, id)
	if err != nil {
		return Get{}, err
	}

	return data, nil
}

func (s Service) GetList(ctx context.Context, filter entity.Filter, authHeader string) ([]Get, int, error) {
	if authHeader != "" {
		isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
		if err != nil {
			return nil, 0, err
		}
		data, count, err := s.repo.GetList(ctx, filter, &isValidToken.Id)
		return data, count, nil
	}

	data, count, err := s.repo.GetList(ctx, filter, nil)
	if err != nil {
		return nil, 0, err
	}
	return data, count, nil
}

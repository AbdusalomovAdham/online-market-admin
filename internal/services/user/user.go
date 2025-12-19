package user

import (
	"context"
	"fmt"
	"main/internal/entity"
	"mime/multipart"
	"time"
)

type Service struct {
	repo Repository
	auth Auth
	file File
}

func NewService(repo Repository, auth Auth, file File) Service {
	return Service{
		repo: repo,
		auth: auth,
		file: file,
	}
}

func (s Service) AdminCreateUser(ctx context.Context, data Create, authHeader string) (int64, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	if data.Password != nil {
		hashPassword, err := s.auth.HashPassword(*data.Password)
		if err != nil {
			return 0, err
		}
		data.Password = &hashPassword
	} else {
		data.Password = nil
	}

	var birthTime *time.Time
	if data.BirthDate != nil {
		layout := "2006-01-02"
		parsedTime, err := time.Parse(layout, *data.BirthDate)
		if err != nil {
			return 0, fmt.Errorf("invalid birth date format: %w", err)
		}
		birthTime = &parsedTime
	}

	return s.repo.Create(ctx, data, isValidToken.Id, *birthTime)
}

func (s Service) MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string, startID *int32) ([]entity.File, error) {
	return s.file.MultipleUpload(ctx, files, folder, startID)
}

func (s Service) Upload(ctx context.Context, image *multipart.FileHeader, folder string) (entity.File, error) {
	return s.file.Upload(ctx, image, folder)
}

func (s Service) GetAll(ctx context.Context, filter entity.Filter) ([]Get, int, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s Service) GetById(ctx context.Context, id int64) (Get, error) {
	return s.repo.GetById(ctx, id)
}

func (s Service) AdminUserUpdate(ctx context.Context, id int64, data Update, authHeader string) error {
	detail, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	isValid, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	if data.Password != nil {
		hashPassword, err := s.auth.HashPassword(*data.Password)
		if err != nil {
			return err
		}
		data.Password = &hashPassword
	} else {
		data.Password = nil
	}

	if detail.Avatar != nil && data.Avatar != nil {
		if err := s.file.Delete(ctx, *detail.Avatar); err != nil {
			return err
		}
	}

	return s.repo.Update(ctx, id, data, isValid.Id)
}

func (s Service) AdminUserDelete(ctx context.Context, id int64, authHeader string) error {
	isValid, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id, isValid.Id)
}

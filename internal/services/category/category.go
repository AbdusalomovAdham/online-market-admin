package category

import (
	"context"
	"main/internal/entity"
)

type Service struct {
	repo Repository
	auth Auth
}

func NewService(repo Repository, auth Auth) Service {
	return Service{repo: repo, auth: auth}
}

func (s *Service) AdminCategoryCreate(ctx context.Context, category Create, authHeader string) (int64, error) {
	token, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, category, token.Id)
}

func (uc *Service) AdminCategoryDelete(ctx context.Context, categoryId int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, categoryId, token.Id)
}

func (uc *Service) AdminCategoryGetById(ctx context.Context, id int64) (CategoryById, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *Service) AdminCategoryGetList(ctx context.Context, filter entity.Filter) ([]Get, int, error) {
	return uc.repo.GetList(ctx, filter)
}

func (uc *Service) AdminCategoryUpdate(ctx context.Context, id int64, data Update, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return uc.repo.Update(ctx, id, data, token.Id)
}

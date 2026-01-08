package param

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

func (s Service) CreateParam(ctx context.Context, param Create, authHeader string) (int64, error) {
	isToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, param, isToken.Id)
}

func (s Service) GetParamById(ctx context.Context, paramId int64) (ParamById, error) {
	return s.repo.GetById(ctx, paramId)
}

func (s *Service) ParamGetList(ctx context.Context, filter entity.Filter) ([]Get, int, error) {
	return s.repo.GetList(ctx, filter)
}

func (s *Service) DeleteParam(ctx context.Context, paramId int64, authHeader string) error {
	isToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, paramId, isToken.Id)
}

func (s *Service) UpdateParam(ctx context.Context, paramId int64, data UpdateParam, authHeader string) error {
	isToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, paramId, data, isToken.Id)
}

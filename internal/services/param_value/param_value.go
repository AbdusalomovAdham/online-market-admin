package param_value

import (
	"context"
	"main/internal/entity"
)

type Service struct {
	repo Repository
	auth Auth
}

func NewService(repo Repository, auth Auth) *Service {
	return &Service{repo: repo, auth: auth}
}

func (s *Service) ParamValueCreate(ctx context.Context, param Create, authHeader string) (int64, error) {
	isToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, param, isToken.Id)
}

func (s *Service) ParamValueGetList(ctx context.Context, filter entity.Filter) ([]Get, int, error) {
	return s.repo.GetList(ctx, filter)
}

func (s *Service) ParamValueDelete(ctx context.Context, paramValueId int64, authHeader string) error {
	isToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, paramValueId, isToken.Id)
}

func (s *Service) ParamValueGetById(ctx context.Context, paramId int64) (ParamValueById, error) {
	return s.repo.GetById(ctx, paramId)
}

func (s *Service) ParamValueUpdate(ctx context.Context, id int64, data Update, authHeader string) error {
	isToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, id, data, isToken.Id)
}

func (s *Service) ParamValueGetListByParamId(ctx context.Context, filter entity.Filter, paramId int) ([]Get, int64, error) {
	return s.repo.GetListByParamId(ctx, filter, paramId)
}

func (s *Service) ParamValueGetByParamId(ctx context.Context, paramId int, filter entity.Filter) ([]ParamValueByParamId, int64, error) {
	return s.repo.GetByParamId(ctx, paramId, filter)
}

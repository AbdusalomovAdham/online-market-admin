package district

import (
	"context"
	"main/internal/entity"
)

type Service struct {
	repo Repository
	auth Auth
}

func NewService(repo Repository, auth Auth) *Service {
	return &Service{
		repo: repo,
		auth: auth,
	}
}

func (s *Service) GetDistricts(ctx context.Context, filter entity.Filter) ([]GetList, int, error) {
	return s.repo.Get(ctx, filter)
}

func (s *Service) GetDistrictsByRegionId(ctx context.Context, filter entity.Filter, regionId int) ([]GetList, error) {
	return s.repo.GetListByRegionId(ctx, filter, regionId)
}

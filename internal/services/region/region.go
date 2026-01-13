package region

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

func (s *Service) GetRegions(ctx context.Context, filter entity.Filter) ([]GetList, int, error) {
	return s.repo.Get(ctx, filter)
}

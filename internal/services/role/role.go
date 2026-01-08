package role

import (
	"context"
	"main/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s *Service) AdminRoleList(ctx context.Context, filter entity.Filter) ([]Get, int, error) {
	return s.repo.GetList(ctx, filter)
}

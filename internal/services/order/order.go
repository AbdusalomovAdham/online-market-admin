package order

import (
	"context"
	"main/internal/entity"
)

type Service struct {
	repo Repository
	auth Auth
}

func NewService(repo Repository, auth Auth) Service {
	return Service{
		repo: repo,
		auth: auth,
	}
}

func (s *Service) AdminOrderCreate(ctx context.Context, order Create, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, order, isValidToken.Id)
}

func (s *Service) AdminOrderGetList(ctx context.Context, authHeader string, filter entity.Filter) ([]Get, int64, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.GetList(ctx, isValidToken.Id, filter)
}

func (s *Service) AdminOrderGetById(ctx context.Context, orderId int64, authHeader string) (Get, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return Get{}, err
	}

	return s.repo.GetById(ctx, orderId, isValidToken.Id)
}

func (s *Service) AdminOrderDelete(ctx context.Context, orderId int64, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, orderId, isValidToken.Id)
}

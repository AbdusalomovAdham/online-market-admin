package wishlist

import (
	"context"
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

func (s *Service) GetList(ctx context.Context, authHeader string) ([]GetList, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return nil, err
	}

	list, err := s.repo.GetList(ctx, isValidToken.Id)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Service) Create(ctx context.Context, productId Create, authHeader string) (int64, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, productId, isValidToken.Id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) Delete(ctx context.Context, wishlistItemId int64, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, wishlistItemId, isValidToken.Id)
}

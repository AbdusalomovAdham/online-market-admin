package wishlist

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

func (s *Service) AdminWishistGetList(ctx context.Context, authHeader string, filter entity.Filter) ([]GetList, int64, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return nil, 0, err
	}

	list, count, err := s.repo.GetList(ctx, isValidToken.Id, filter)
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (s *Service) AdminWishlistCreate(ctx context.Context, productId Create, authHeader string) (int64, error) {
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

func (s *Service) AdminWishlistDelete(ctx context.Context, wishlistItemId int64, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, wishlistItemId, isValidToken.Id)
}

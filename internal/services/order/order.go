package order

import "context"

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

func (s *Service) Create(ctx context.Context, order Create, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, order, isValidToken.Id)
}

func (s *Service) GetList(ctx context.Context, authHeader string) ([]Get, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return nil, err
	}

	return s.repo.GetList(ctx, isValidToken.Id)
}

func (s *Service) GetById(ctx context.Context, orderId int64, authHeader string) (Get, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return Get{}, err
	}

	return s.repo.GetById(ctx, orderId, isValidToken.Id)
}

func (s *Service) Delete(ctx context.Context, orderId int64, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, orderId, isValidToken.Id)
}

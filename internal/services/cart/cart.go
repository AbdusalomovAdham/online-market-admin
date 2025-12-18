package cart

import "context"

type Service struct {
	repo Repository
	auth Auth
}

func NewService(repo Repository, auth Auth) Service {
	return Service{repo: repo, auth: auth}
}

func (s Service) Create(ctx context.Context, cart Create, authHeader string) (int64, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, cart.ProductId, isValidToken.Id)
}

func (s Service) UpdateCartItemTotal(ctx context.Context, cartItemId int64, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, cartItemId, isValidToken.Id)
}

func (s Service) DeleteCartItem(ctx context.Context, cartItemId int64, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.DeleteCartItem(ctx, cartItemId, isValidToken.Id)
}

func (s Service) GetList(ctx context.Context, authHeader string) ([]Get, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return nil, err
	}

	return s.repo.GetList(ctx, isValidToken.Id)
}

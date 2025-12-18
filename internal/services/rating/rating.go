package rating

import (
	"context"
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

func (s Service) CreateRating(ctx context.Context, data Create, authHeader string) (int64, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, data, isValidToken.Id)
}

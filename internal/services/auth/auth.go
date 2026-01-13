package auth

import (
	"context"
	"fmt"
	"main/internal/usecase/auth"
)

type Service struct {
	repo    Repository
	auth    Auth
	cache   Cache
	sendSMS SendSMS
}

func NewService(repo Repository, cache Cache, sendSMS SendSMS, auth Auth) *Service {
	return &Service{
		repo:    repo,
		auth:    auth,
		cache:   cache,
		sendSMS: sendSMS,
	}
}

func (s *Service) SignIn(ctx context.Context, data SignIn) (AdminDetails, string, error) {
	detail, err := s.repo.GetByLogin(ctx, data.Login)
	if err != nil {
		return AdminDetails{}, "", err
	}

	if detail.Password == "" || !s.auth.CheckPasswordHash(data.Password, detail.Password) {
		return AdminDetails{}, "", fmt.Errorf("Error password or login")
	}

	var generateToken auth.GenerateToken
	generateToken.Id = detail.Id
	generateToken.Role = detail.Role

	token, err := s.auth.GenerateToken(ctx, generateToken)
	if err != nil {
		return AdminDetails{}, "", err
	}

	return detail, token, nil
}

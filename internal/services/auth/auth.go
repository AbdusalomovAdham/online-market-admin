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

func (s *Service) SignIn(ctx context.Context, data SignIn) (string, error) {
	password, userId, roleId, err := s.repo.GetByLogin(ctx, data.Login)
	if err != nil {
		return "", err
	}

	if password == "" || !s.auth.CheckPasswordHash(data.Password, password) {
		return "", fmt.Errorf("Error password or login")
	}

	var generateToken auth.GenerateToken
	generateToken.Id = userId
	generateToken.Role = roleId

	token, err := s.auth.GenerateToken(ctx, generateToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

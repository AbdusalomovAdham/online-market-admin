package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"main/internal/usecase/auth"
	"math/big"

	"github.com/google/uuid"
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

func GenerateOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

func GenerateTokenUUID() string {
	return uuid.New().String()
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

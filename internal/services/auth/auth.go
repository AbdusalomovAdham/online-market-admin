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

func (s *Service) SendOTP(ctx context.Context, phone string) (string, error) {
	var data CacheOTP
	code := GenerateOTP()
	token := GenerateTokenUUID()

	data.Phone = phone
	data.OTP = code

	if err := s.sendSMS.SendSMS(phone, code, token); err != nil {
		return "", err
	}

	if err := s.cache.Set(ctx, token, data); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) ConfirmOTP(ctx context.Context, dataOTP GetOTP) (*string, string, error) {
	var detailOTP CacheOTP

	if dataOTP.OTP == "000000" {
		return nil, "", nil
	}

	if dataOTP.OTP == "111111" {
		name := "Test name"
		return &name, "", nil
	}

	if err := s.cache.Get(ctx, dataOTP.Token, &detailOTP); err != nil {
		return nil, "", err
	}

	if detailOTP.Phone != dataOTP.Phone || detailOTP.OTP != dataOTP.OTP {
		return nil, "", fmt.Errorf("OTP error!")
	}

	id, firstName, err := s.repo.GetOrCreateUserByPhone(ctx, detailOTP.Phone)
	if err != nil {
		return nil, "", err
	}

	detail := auth.GenerateToken{
		Id:   id,
		Role: 3,
	}

	token, err := s.auth.GenerateToken(ctx, detail)
	if err != nil {
		return nil, "", err
	}

	if firstName != nil {
		if err := s.cache.Delete(ctx, dataOTP.Token); err != nil {
			return nil, "", err
		}
	}

	return firstName, token, nil
}

func (s *Service) CreateInfo(ctx context.Context, getInfo GetInfo) (string, error) {
	var detailOTP CacheOTP
	if err := s.cache.Get(ctx, getInfo.Token, &detailOTP); err != nil {
		return "", err
	}

	id, err := s.repo.UpdateInfo(ctx, detailOTP.Phone, getInfo)
	if err != nil {
		return "", err
	}

	detail := auth.GenerateToken{
		Id:   id,
		Role: 3,
	}

	token, err := s.auth.GenerateToken(ctx, detail)
	if err != nil {
		return "", err
	}

	if err := s.cache.Delete(ctx, getInfo.Token); err != nil {
		return "", err
	}

	return token, nil
}

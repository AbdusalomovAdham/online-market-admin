package auth

import (
	"context"
	"main/internal/entity"
	"main/internal/usecase/auth"
)

type Repository interface {
	GetByLogin(ctx context.Context, login string) (AdminDetails, error)
}

type Auth interface {
	GenerateToken(ctx context.Context, data auth.GenerateToken) (string, error)
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateResetToken(n int) (string, error)
}

type SendSMS interface {
	SendSMS(phone, code, token string) error
}

type Email interface {
	SendEmail(email, subject, message string) error
}

type Cache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, key string) error
}

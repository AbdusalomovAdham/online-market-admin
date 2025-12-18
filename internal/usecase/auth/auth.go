package auth

import (
	"context"
	"encoding/hex"
	"errors"
	"main/internal/entity"
	"main/internal/pkg/config"

	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GenerateToken(ctx context.Context, data GenerateToken) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["role"] = data.Role
	claims["id"] = data.Id

	tokenStr, err := token.SignedString([]byte(config.GetConfig().JWTKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (uc *UseCase) IsValidToken(ctx context.Context, tokenStr string) (entity.User, error) {
	claims := new(Claims)
	token := strings.TrimPrefix(tokenStr, "Bearer ")
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.GetConfig().JWTKey), nil
	})

	if err != nil {
		return entity.User{}, err
	}

	if !tkn.Valid {
		return entity.User{}, errors.New("invalid token")
	}

	userDetail, err := uc.repo.GetById(ctx, claims.ID)
	if err != nil {
		return entity.User{}, errors.New("user not found")
	}
	return userDetail, nil
}

func (uc *UseCase) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	return string(bytes), err
}

func (uc *UseCase) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (uc *UseCase) GenerateResetToken(n int) (string, error) {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

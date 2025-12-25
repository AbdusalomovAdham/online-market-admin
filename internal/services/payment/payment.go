package payment

import (
	"context"
	"main/internal/entity"
)

type Service struct {
	repo Repository
	auth Auth
}

func NewService(repo Repository, auth Auth) Service {
	return Service{repo: repo, auth: auth}
}

func (s *Service) AdminPaymentCreate(ctx context.Context, paymentStatus Create, authHeader string) (int64, error) {
	token, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, paymentStatus, token.Id)
}

func (uc *Service) AdminPaymentDelete(ctx context.Context, paymentStatusId int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, paymentStatusId, token.Id)
}

func (uc *Service) AdminPaymentGetById(ctx context.Context, id int64) (PaymentStatusById, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *Service) AdminPaymentGetList(ctx context.Context, filter entity.Filter, lang string) ([]Get, int, error) {
	return uc.repo.GetList(ctx, filter, lang)
}

func (uc *Service) AmdinPaymentUpdate(ctx context.Context, id int64, data Update, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return uc.repo.Update(ctx, id, data, token.Id)
}

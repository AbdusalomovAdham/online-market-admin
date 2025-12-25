package order_status

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

func (s *Service) AdminOrderStatusCreate(ctx context.Context, orderStatus Create, authHeader string) (int64, error) {
	token, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, orderStatus, token.Id)
}

func (uc *Service) AdminOrderStatusDelete(ctx context.Context, orderStatusId int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, orderStatusId, token.Id)
}

func (uc *Service) AdminOrderStatusGetById(ctx context.Context, id int64) (OrderStatusById, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *Service) AdminOrderStatusGetList(ctx context.Context, filter entity.Filter) ([]Get, int, error) {
	return uc.repo.GetList(ctx, filter)
}

func (uc *Service) AdminOrderStatusUpdate(ctx context.Context, id int64, data Update, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return uc.repo.Update(ctx, id, data, token.Id)
}

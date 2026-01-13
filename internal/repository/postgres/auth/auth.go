package auth

import (
	"context"
	"main/internal/entity"
	"main/internal/services/auth"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) GetByLogin(ctx context.Context, login string) (auth.AdminDetails, error) {
	var details auth.AdminDetails
	query := `SELECT id, password, role, first_name, last_name, avatar FROM users WHERE deleted_at IS NULL AND login = ? AND status = true`

	err := r.QueryRowContext(ctx, query, login).Scan(&details.Id, &details.Password, &details.Role, &details.FirstName, &details.LastName, &details.Avatar)
	if err != nil {
		return auth.AdminDetails{}, err
	}

	return details, nil
}

func (r Repository) GetById(ctx context.Context, id int) (entity.User, error) {
	var detail entity.User
	query := `SELECT id, first_name, last_name, phone_number FROM users WHERE id = ?`
	err := r.QueryRowContext(ctx, query, id).Scan(&detail.Id, &detail.FirstName, &detail.LastName, &detail.PhoneNumber)
	if err != nil {
		return entity.User{}, err
	}
	return detail, nil
}

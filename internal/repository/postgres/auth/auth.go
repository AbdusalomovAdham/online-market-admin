package auth

import (
	"context"
	"database/sql"
	"main/internal/entity"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) GetByLogin(ctx context.Context, login string) (string, int64, int, error) {
	var id int64
	var role int
	var password sql.NullString

	query := `SELECT id, password, role FROM users WHERE deleted_at IS NULL AND login = ? AND status = true`

	err := r.QueryRowContext(ctx, query, login).Scan(&id, &password, &role)
	if err != nil {
		return "", 0, 0, err
	}

	pw := ""
	if password.Valid {
		pw = password.String
	}

	return pw, id, role, nil
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

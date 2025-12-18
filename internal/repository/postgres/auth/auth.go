package auth

import (
	"context"
	"database/sql"
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

func (r Repository) Create(ctx context.Context, user *auth.Create, phoneNumber string) error {
	var id int64
	query := `INSERT INTO users (first_name, last_name, created_at) VALUES (?, ?, ?, NOW()) WHERE phone_number = ? RETURNING id`

	err := r.QueryRowContext(ctx, query, user.FirstName, user.LastName, phoneNumber).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetOrCreateUserByPhone(ctx context.Context, phone string) (int64, *string, error) {
	var firstName string
	var id int64
	query := `SELECT id, first_name FROM users WHERE phone_number = ?`
	err := r.QueryRowContext(ctx, query, phone).Scan(&id, &firstName)
	if err == nil {
		return id, &firstName, nil
	}

	if err != sql.ErrNoRows {
		return 0, nil, err
	} else {
		insertQuery := `INSERT INTO users (phone_number, role, created_at) VALUES (?, ?, NOW()) RETURNING id, first_name`
		err := r.QueryRowContext(ctx, insertQuery, phone, 3).Scan(&id, &firstName)
		if err != nil {
			return 0, nil, err
		}
	}
	return id, nil, nil
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

func (r Repository) UpdateInfo(ctx context.Context, phoneNumber string, info auth.GetInfo) (int64, error) {
	var id int64
	query := `UPDATE users SET first_name = ?, last_name = ? WHERE phone_number = ? RETURNING id`
	err := r.QueryRowContext(ctx, query, info.FirstName, info.LastName, phoneNumber).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

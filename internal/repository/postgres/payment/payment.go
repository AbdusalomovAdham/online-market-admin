package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/entity"
	"main/internal/services/payment"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, paymentStatus payment.Create, adminId int64) (int64, error) {
	var id int64

	query := `INSERT INTO payments (name, created_by, status) VALUES (?, ?, COALESCE(?, TRUE)) RETURNING id`
	if err := r.QueryRowContext(ctx, query, paymentStatus.Name, adminId, paymentStatus.Status).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) Delete(ctx context.Context, paymenStatusId int64, adminId int64) error {

	query := `UPDATE payments
			SET deleted_at = NOW(), deleted_by = ?
			WHERE id = ? RETURNING id`
	_, err := r.ExecContext(ctx, query, adminId, paymenStatusId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, paymentStatusid int64) (payment.PaymentStatusById, error) {

	var detail payment.PaymentStatusById

	query := `
        SELECT id, status, created_at, name
        FROM payments
        WHERE id = ? AND deleted_at IS NULL
    `

	err := r.QueryRowContext(ctx, query, paymentStatusid).
		Scan(&detail.Id, &detail.Status, &detail.CreatedAt, &detail.Name)

	if err != nil {
		return payment.PaymentStatusById{}, err
	}

	return detail, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter, lang string) ([]payment.Get, int, error) {
	if lang == "" {
		lang = "uz"
	}

	var list []payment.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE os.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY os.id DESC"
	if filter.Order != nil && *filter.Order != "" {
		parts := strings.Fields(*filter.Order)
		if len(parts) == 2 {
			column := parts[0]
			direction := strings.ToUpper(parts[1])
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC"
			}
			orderQuery = fmt.Sprintf("ORDER BY %s %s", column, direction)
		}
	}

	query := fmt.Sprintf(`
	    SELECT
	        os.id,
	        os.name->>'%s' as name,
	        os.status,
	        os.created_at
	    FROM payments os
	    %s %s %s %s
	`, lang, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &list); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(os.id)
				FROM payments os
				WHERE os.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select payment status count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, data payment.Update, userId int64) error {
	var paymentStatus entity.Payment
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by FROM payments WHERE id = ? AND deleted_at is NULL`
	if err := r.QueryRowContext(ctx, query, id).Scan(&paymentStatus.Id, &nameJSON, &paymentStatus.Status, &paymentStatus.UpdatedAt, &paymentStatus.UpdatedBy); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &paymentStatus.Name); err != nil {
		return err
	}

	if data.Name != nil {
		if data.Name.Uz != nil {
			paymentStatus.Name.Uz = data.Name.Uz
		}

		if data.Name.Ru != nil {
			paymentStatus.Name.Ru = data.Name.Ru
		}

		if data.Name.En != nil {
			paymentStatus.Name.En = data.Name.En
		}
	}

	if data.Status != nil {
		paymentStatus.Status = data.Status
	}

	query = `UPDATE payments SET name = ?, status = ?, updated_by = ?, updated_at = NOW() WHERE id = ? AND deleted_at is NULL`
	if _, err := r.ExecContext(ctx, query, paymentStatus.Name, paymentStatus.Status, userId, id); err != nil {
		return err
	}
	return nil
}

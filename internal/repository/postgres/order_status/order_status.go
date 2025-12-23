package order_status

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/entity"
	"main/internal/services/order_status"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, category order_status.Create, userId int64) (int64, error) {
	var id int64

	query := `INSERT INTO order_statuses (name, created_by, status) VALUES (?, ?, COALESCE(?, TRUE)) RETURNING id`
	if err := r.QueryRowContext(ctx, query, category.Name, userId, category.Status).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) Delete(ctx context.Context, id int64, userId int64) error {

	query := `UPDATE order_statuses SET deleted_at = NOW(), deleted_by = ? WHERE id = ?  RETURNING id`
	_, err := r.ExecContext(ctx, query, userId, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (order_status.OrderStatusById, error) {

	var detail order_status.OrderStatusById

	query := `
        SELECT id, status, created_at, name
        FROM order_statuses
        WHERE id = ? AND deleted_at IS NULL
    `

	err := r.QueryRowContext(ctx, query, id).
		Scan(&detail.Id, &detail.Status, &detail.CreatedAt, &detail.Name)

	if err != nil {
		return order_status.OrderStatusById{}, err
	}

	return detail, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter, lang string) ([]order_status.Get, int, error) {
	if lang == "" {
		lang = "uz"
	}

	var list []order_status.Get
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
	    FROM order_statuses os
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

	countQuery := `SELECT COUNT(os.id) FROM order_statuses os WHERE os.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select order status count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, data order_status.Update, userId int64) error {
	var orderStatus entity.OrderStatus
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by FROM order_statuses WHERE id = ? AND deleted_at is NULL`
	if err := r.QueryRowContext(ctx, query, id).Scan(&orderStatus.Id, &nameJSON, &orderStatus.Status, &orderStatus.UpdatedAt, &orderStatus.UpdatedBy); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &orderStatus.Name); err != nil {
		return err
	}

	if data.Name != nil {
		if data.Name.Uz != nil {
			orderStatus.Name.Uz = data.Name.Uz
		}

		if data.Name.Ru != nil {
			orderStatus.Name.Ru = data.Name.Ru
		}

		if data.Name.En != nil {
			orderStatus.Name.En = data.Name.En
		}
	}

	if data.Status != nil {
		orderStatus.Status = data.Status
	}

	query = `UPDATE order_statuses SET name = ?, status = ?, updated_by = ?, updated_at = NOW() WHERE id = ? AND deleted_at is NULL`
	if _, err := r.ExecContext(ctx, query, orderStatus.Name, orderStatus.Status, userId, id); err != nil {
		return err
	}
	return nil
}

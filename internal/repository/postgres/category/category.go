package category

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/entity"
	"main/internal/services/category"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, category category.Create, userId int64) (int64, error) {
	var id int64

	query := `INSERT INTO categories (name, created_by, status, parent_id) VALUES (?, ?, COALESCE(?, TRUE), ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, category.Name, userId, category.Status, category.ParentId).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) Delete(ctx context.Context, id int64, userId int64) error {

	query := `UPDATE categories SET deleted_at = NOW(), deleted_by = ? WHERE id = ?  RETURNING id`

	_, err := r.ExecContext(ctx, query, userId, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (category.CategoryById, error) {
	var detail category.CategoryById

	query := `SELECT id, status, created_at, name FROM categories WHERE id = ? AND deleted_at is NULL`

	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return category.CategoryById{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &detail); err != nil {
		return category.CategoryById{}, err
	}
	return detail, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter, lang string) ([]category.Get, int, error) {
	var list []category.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE d.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY d.id DESC"
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
	        d.id,
	        d.name->>'%s' as name,
	        d.status,
	        d.created_at
	    FROM categories d
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

	countQuery := `SELECT COUNT(d.id) FROM categories d WHERE d.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0
	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select category count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Update(ctx context.Context, id int64, data category.Update, userId int64) error {
	var category entity.Category
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by FROM categories WHERE id = ? AND deleted_at is NULL`
	if err := r.QueryRowContext(ctx, query, id).Scan(&category.Id, &nameJSON, &category.Status, &category.UpdatedAt, &category.UpdatedBy); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &category.Name); err != nil {
		return err
	}

	if data.Name != nil {
		if data.Name.Uz != nil {
			category.Name.Uz = data.Name.Uz
		}

		if data.Name.Ru != nil {
			category.Name.Ru = data.Name.Ru
		}

		if data.Name.En != nil {
			category.Name.En = data.Name.En
		}
	}

	if data.Status != nil {
		category.Status = data.Status
	}

	query = `UPDATE categories SET name = ?, status = ?, updated_by = ?, updated_at = NOW() WHERE id = ? AND deleted_at is NULL`
	if _, err := r.ExecContext(ctx, query, category.Name, category.Status, userId, id); err != nil {
		return err
	}
	return nil
}

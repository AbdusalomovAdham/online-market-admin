package param

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/entity"
	"main/internal/services/category"
	"main/internal/services/param"
	"strings"

	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, param param.Create, adminId int64) (int64, error) {
	var newParamId int64

	query := `INSERT INTO params (name, status, created_by, created_at, type) VALUES (?, COALESCE(?, TRUE), ?, NOW(), ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, param.Name, *param.Status, adminId, param.Type).Scan(&newParamId); err != nil {
		return 0, err
	}

	if newParamId > 0 {
		query = `INSERT INTO category_params (category_id, param_id, created_by, created_at) VALUES (?, ?, ?, NOW())`
		if _, err := r.ExecContext(ctx, query, pq.Array(param.CategoryId), newParamId, adminId); err != nil {
			return 0, err
		}
	}

	return newParamId, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter) ([]param.Get, int, error) {
	var list []param.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE cp.deleted_at IS NULL AND p.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY p.id DESC"
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

	if filter.Search != nil {
		searchQuery := fmt.Sprintf(" AND (p.name ->> '%s') ILIKE '%%%s%%'", *filter.Language, *filter.Search)
		whereQuery += searchQuery
	}

	query := fmt.Sprintf(`
	    SELECT
	        p.id,
	        p.status,
	        p.created_at,
	        p.name->>'%s' as param_name,
			p.type as type
	    FROM category_params cp
		LEFT JOIN params p ON cp.param_id = p.id
	    %s %s %s %s`, *filter.Language, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &list); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(p.id) FROM params p WHERE p.deleted_at IS NULL AND p.id IN (SELECT param_id FROM category_params WHERE deleted_at IS NULL)`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0
	if len(list) == 0 {
		return list, count, nil
	}

	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select param count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Delete(ctx context.Context, paramId int64, adminId int64) error {

	query := `UPDATE params SET deleted_at = NOW(), deleted_by = ? WHERE id = ?  RETURNING id`
	_, err := r.ExecContext(ctx, query, adminId, paramId)
	if err != nil {
		return err
	}

	query = `UPDATE category_params SET deleted_at = NOW(), deleted_by = ? WHERE param_id = ?`
	_, err = r.ExecContext(ctx, query, adminId, paramId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetById(ctx context.Context, paramId int64) (param.ParamById, error) {
	var detail param.ParamById

	query := `
	    SELECT
	        p.id,
	        p.status,
	        p.created_at,
			p.type,
	        p.name,
			cp.category_id
	    FROM params p
		LEFT JOIN category_params cp ON p.id = cp.param_id
		WHERE p.id = ? AND p.deleted_at IS NULL`

	row := r.QueryRowContext(ctx, query, paramId)
	var rawName []byte
	if err := row.Scan(
		&detail.Id,
		&detail.Status,
		&detail.CreatedAt,
		&detail.Type,
		&rawName,
		pq.Array(&detail.CategoryId),
	); err != nil {
		return param.ParamById{}, err
	}

	json.Unmarshal(rawName, &detail.Name)

	return detail, nil
}

func (r *Repository) Update(ctx context.Context, paramId int64, data param.UpdateParam, adminId int64) error {
	var param entity.Param
	var oldCategoryId []int64
	var nameJSON []byte

	query := `SELECT id, name, status, updated_at, updated_by, type FROM params WHERE id = ? AND deleted_at is NULL`
	if err := r.QueryRowContext(ctx, query, paramId).Scan(&param.Id, &nameJSON, &param.Status, &param.UpdatedAt, &param.UpdatedBy, &param.Type); err != nil {
		return err
	}

	query = `SELECT category_id FROM category_params WHERE param_id = ?`
	if err := r.QueryRowContext(ctx, query, paramId).Scan(pq.Array(&oldCategoryId)); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &param.Name); err != nil {
		return err
	}

	if data.Name != nil {
		if data.Name.Uz != nil {
			param.Name.Uz = data.Name.Uz
		}

		if data.Name.Ru != nil {
			param.Name.Ru = data.Name.Ru
		}

		if data.Name.En != nil {
			param.Name.En = data.Name.En
		}
	}

	if data.Status != nil {
		param.Status = *data.Status
	}

	if data.CategoryId != nil {
		query := `UPDATE category_params SET category_id = ? WHERE param_id = ?`
		if _, err := r.ExecContext(ctx, query, pq.Array(data.CategoryId), paramId); err != nil {
			return err
		}
	}

	if data.Type != nil {
		param.Type = *data.Type
	}

	query = `UPDATE params SET name = ?, status = ?, type = ?, updated_by = ?, updated_at = NOW() WHERE id = ? AND deleted_at is NULL`
	if _, err := r.ExecContext(ctx, query, param.Name, param.Status, param.Type, adminId, paramId); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByParentId(ctx context.Context, filter entity.Filter, categoryParentId int64) ([]category.Get, int, error) {
	var list []category.Get
	var limitQuery, offsetQuery string
	log.Println("asdf", categoryParentId)
	whereQuery := fmt.Sprintf("WHERE c.deleted_at IS NULL AND c.parent_id = %d", categoryParentId)

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY c.id DESC"
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
	        c.id,
	        c.name->>'%s' as name,
	        c.status,
	        c.created_at
	    FROM categories c
	    %s %s %s %s
	`, *filter.Language, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &list); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(c.id) FROM categories c WHERE c.deleted_at IS NULL`

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

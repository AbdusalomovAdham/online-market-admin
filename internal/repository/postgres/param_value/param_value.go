package param_value

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"main/internal/entity"
	"main/internal/services/param_value"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, param param_value.Create, adminId int64) (int64, error) {
	var id int64

	query := `INSERT INTO param_values (name, param_id, status, created_by, created_at) VALUES (?, ?, COALESCE(?, TRUE), ?, NOW()) RETURNING id`
	if err := r.QueryRowContext(ctx, query, param.Name, param.ParamId, param.Status, adminId).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) Delete(ctx context.Context, paramValueId int64, adminId int64) error {

	query := `UPDATE param_values SET deleted_at = NOW(), deleted_by = ? WHERE id = ?  RETURNING id`

	_, err := r.ExecContext(ctx, query, adminId, paramValueId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, paramId int64) (param_value.ParamValueById, error) {
	var detail param_value.ParamValueById
	var paramValeuName []byte
	query := `
		SELECT
			pv.id,
			pv.status,
			pv.created_at,
			pv.name,
			pv.param_id
		FROM param_values pv
		LEFT JOIN params p
		  ON p.id = pv.param_id
		 AND p.deleted_at IS NULL
		WHERE pv.id = ?
		  AND pv.deleted_at IS NULL
	`

	err := r.QueryRowContext(ctx, query, paramId).Scan(
		&detail.Id,
		&detail.Status,
		&detail.CreatedAt,
		&paramValeuName,
		&detail.ParamId,
	)

	if err != nil {
		if err == sql.ErrNoRows {

			return param_value.ParamValueById{}, sql.ErrNoRows
		}
		return param_value.ParamValueById{}, err
	}

	json.Unmarshal(paramValeuName, &detail.Name)

	return detail, nil
}

func (r *Repository) GetList(ctx context.Context, filter entity.Filter) ([]param_value.Get, int, error) {
	var list []param_value.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE pv.deleted_at IS NULL"

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY pv.id DESC"
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
		searchQuery := fmt.Sprintf(" AND (pv.name ->> '%s') ILIKE '%%%s%%'", *filter.Language, *filter.Search)
		whereQuery += searchQuery
	}

	query := fmt.Sprintf(`
	    SELECT
	        pv.id,
	        pv.name->>'%s' as name,
	        pv.status,
	        pv.created_at,
			pv.param_id,
			p.name->>'%s' as param_name
	    FROM param_values pv
		LEFT JOIN params p ON pv.param_id = p.id
	    %s %s %s %s
	`, *filter.Language, *filter.Language, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &list); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(pv.id) FROM param_values pv WHERE pv.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0
	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select param value count: %w", err)
	}
	return list, count, nil
}

func (r *Repository) Update(ctx context.Context, paramId int64, data param_value.Update, adminId int64) error {
	var paramValue entity.ParamValue
	var nameJSON []byte
	query := `SELECT id, name, status, updated_at, updated_by FROM param_values WHERE id = ? AND deleted_at is NULL`
	if err := r.QueryRowContext(ctx, query, paramId).Scan(&paramValue.Id, &nameJSON, &paramValue.Status, &paramValue.UpdatedAt, &paramValue.UpdatedBy); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &paramValue.Value); err != nil {
		return err
	}

	if data.Name != nil {
		if data.Name.Uz != nil {
			paramValue.Value.Uz = data.Name.Uz
		}

		if data.Name.Ru != nil {
			paramValue.Value.Ru = data.Name.Ru
		}

		if data.Name.En != nil {
			paramValue.Value.En = data.Name.En
		}
	}

	if data.Status != nil {
		paramValue.Status = *data.Status
	}

	if data.ParamId != nil {
		paramValue.ParamId = *data.ParamId
	}
	query = `UPDATE param_values SET name = ?, status = ?, updated_by = ?, updated_at = NOW() ,param_id = ? WHERE id = ? AND deleted_at is NULL`
	if _, err := r.ExecContext(ctx, query, paramValue.Value, paramValue.Status, adminId, paramValue.ParamId, paramId); err != nil {
		return err
	}
	return nil
}

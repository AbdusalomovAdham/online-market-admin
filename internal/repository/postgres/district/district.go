package district

import (
	"context"
	"fmt"
	"main/internal/entity"
	"main/internal/services/district"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Get(ctx context.Context, filter entity.Filter) ([]district.GetList, int, error) {
	var districtsList []district.GetList
	var limitQuery, offsetQuery string

	whereQuery := `WHERE deleted_at IS NULL`

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := `ORDER BY id ASC`
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
		searchQuery := fmt.Sprintf(" AND (c.name ->> '%s' ILIKE '%%%s%%')", *filter.Language, *filter.Search)
		whereQuery += searchQuery
	}

	query := fmt.Sprintf(
		`SELECT
			id,
			name ->> '%s' AS name,
			created_at,
			region_id
		FROM districts
		%s %s %s %s`,
		*filter.Language, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &districtsList); err != nil {
		return nil, 0, err
	}

	count := 0
	countQuery := `SELECT COUNT(id) FROM districts WHERE deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	if err := r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, err
	}

	return districtsList, count, nil
}

func (r *Repository) GetListByRegionId(ctx context.Context, filter entity.Filter, regionId int) ([]district.GetList, error) {
	var districtsList []district.GetList
	query := fmt.Sprintf(
		`SELECT
			id,
			name ->> '%s' AS name,
			created_at,
			region_id
		FROM districts
		WHERE region_id = %d`,
		*filter.Language, regionId)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &districtsList); err != nil {
		return nil, err
	}

	return districtsList, nil
}

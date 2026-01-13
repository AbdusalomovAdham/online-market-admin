package region

import (
	"context"
	"fmt"
	"main/internal/entity"
	"main/internal/services/region"
	"strings"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Get(ctx context.Context, filter entity.Filter) ([]region.GetList, int, error) {
	var regionsList []region.GetList
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
			created_at
		FROM regions
		%s %s %s %s`,
		*filter.Language, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &regionsList); err != nil {
		return nil, 0, err
	}

	count := 0
	countQuery := `SELECT COUNT(id) FROM regions WHERE deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	if err := r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, err
	}

	return regionsList, count, nil
}

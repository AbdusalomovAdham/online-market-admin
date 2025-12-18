package product

import (
	"context"
	"fmt"
	"log"
	"main/internal/entity"
	product "main/internal/services/product"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) Create(ctx context.Context, data product.Create, userId int64) (int64, error) {
	var id int64

	query := `INSERT INTO products (name, description, price, stock_quantity, category_id, discount_percent, images, created_by, status, created_at, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, data.Name, data.Description, data.Price, data.StockQuantity, data.CategoryId, data.DiscountPercent, data.Images, userId, false, time.Now(), userId).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func IncrementViewCount(ctx context.Context, productId int64, r *Repository) error {
	query := `
		UPDATE products
		SET views_count = views_count + 1
		WHERE id = ? AND deleted_at IS NULL AND status = true
	`
	_, err := r.ExecContext(ctx, query, productId)
	return err
}

func (r Repository) GetById(ctx context.Context, id int64) (product.Get, error) {
	IncrementViewCount(ctx, id, &r)
	var data product.Get
	log.Println("id", id)
	query := `
			SELECT
			p.id,
			p.name,
			p.description,
			p.price,
			p.stock_quantity,
			p.rating_avg,
			p.seller_id,
			p.category_id,
			p.views_count,
			p.discount_percent,
			p.images,
			p.created_at,
			u.first_name,
			u.last_name,
			u.avatar
			FROM products p
			LEFT JOIN users u ON p.seller_id = u.id
			WHERE p.id = ? AND p.deleted_at IS NULL AND p.status = true
		`
	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return product.Get{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &data); err != nil {
		return product.Get{}, err
	}

	return data, nil
}

func (r Repository) GetList(ctx context.Context, filter entity.Filter, userId *int64) ([]product.Get, int, error) {
	var data []product.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE p.deleted_at IS NULL AND p.status = true"
	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	if filter.CategoryId != nil {
		whereQuery += fmt.Sprintf(" AND p.category_id = %d", *filter.CategoryId)
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

	query := fmt.Sprintf(`
		SELECT
			p.id,
			p.name,
			p.description,
			p.price,
			p.stock_quantity,
			p.rating_avg,
			p.seller_id,
			p.category_id,
			p.views_count,
			p.discount_percent,
			p.images
		FROM products p
		%s
		%s
		%s
		%s
	`, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &data); err != nil {
		return nil, 0, err
	}

	wishlistMap := make(map[int64]bool)
	if userId != nil {
		wlQuery := `SELECT product_id FROM wishlists WHERE user_id = ? AND status = true AND deleted_at IS NULL`
		wlRows, err := r.QueryContext(ctx, wlQuery, *userId)
		if err == nil {
			defer wlRows.Close()
			var pid int64
			for wlRows.Next() {
				if err := wlRows.Scan(&pid); err == nil {
					wishlistMap[pid] = true
				}
			}
		}
	}

	for i := range data {
		data[i].IsWishlist = wishlistMap[data[i].Id]
	}

	countQuery := `SELECT COUNT(p.id) FROM products p WHERE p.deleted_at IS NULL AND p.status = true`
	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := 0
	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select product count: %w", err)
	}

	return data, count, nil
}

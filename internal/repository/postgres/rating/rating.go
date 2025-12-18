package rating

import (
	"context"
	"fmt"
	"main/internal/services/rating"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) Create(ctx context.Context, data rating.Create, customerId int64) (int64, error) {
	var id int64

	tx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var exists bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM products WHERE id = ? AND deleted_at IS NULL AND status = true)
	`, data.ProductId).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, fmt.Errorf("product with id %d not found", data.ProductId)
	}

	query := `
		INSERT INTO ratings (product_id, customer_id, rating, created_at)
		VALUES (?, ?, ?, NOW())
		ON CONFLICT (customer_id, product_id)
		DO UPDATE SET
			rating = EXCLUDED.rating,
			updated_at = NOW()
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query, data.ProductId, customerId, data.Rating).Scan(&id)
	if err != nil {
		return 0, err
	}

	updateProductQuery := `
		UPDATE products
		SET
			rating_avg = (
				SELECT ROUND(AVG(rating)::numeric,1) FROM ratings WHERE product_id = ? AND deleted_at IS NULL AND status = true
			),
			rating_count = (
				SELECT COUNT(*) FROM ratings WHERE product_id = ? AND deleted_at IS NULL AND status = true
			)
		WHERE id = ?
	`

	_, err = tx.ExecContext(ctx, updateProductQuery, data.ProductId, data.ProductId, data.ProductId)
	if err != nil {
		return 0, err
	}

	return id, nil
}

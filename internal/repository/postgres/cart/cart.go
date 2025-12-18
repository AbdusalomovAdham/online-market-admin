package cart

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"main/internal/entity"
	"main/internal/services/cart"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r Repository) Create(ctx context.Context, productId, customerId int64) (int64, error) {
	var cartId int64

	getCartQuery := `
		SELECT id FROM carts
		WHERE customer_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.QueryRowContext(ctx, getCartQuery, customerId).Scan(&cartId)
	if err != nil {
		if err == sql.ErrNoRows {
			createCartQuery := `
				INSERT INTO carts (
					customer_id,
					status,
					total_amount,
					created_at,
					created_by
				)
				VALUES (?, true, 0, NOW(), ?)
				RETURNING id
			`

			if err := r.QueryRowContext(
				ctx,
				createCartQuery,
				customerId,
				customerId,
			).Scan(&cartId); err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	insertItemQuery := `
		INSERT INTO cart_items (
			cart_id,
			product_id,
			quantity,
			price,
			status,
			created_at,
			created_by
		)
		VALUES (
			?, ?, 1,
			(SELECT price FROM products WHERE id = ?),
			true,
			NOW(),
			?
		)
	`

	if _, err := r.ExecContext(
		ctx,
		insertItemQuery,
		cartId,
		productId,
		productId,
		customerId,
	); err != nil {
		return 0, err
	}

	updateTotalQuery := `
		UPDATE carts
		SET total_amount = (
			SELECT COALESCE(SUM(ci.quantity * ci.price), 0)
			FROM cart_items ci
			WHERE ci.cart_id = ?
			  AND ci.deleted_at IS NULL
		)
		WHERE id = ?
		`

	if _, err := r.ExecContext(ctx, updateTotalQuery, cartId, cartId); err != nil {
		return 0, err
	}

	return cartId, nil
}

func (r Repository) Update(ctx context.Context, cartItemId int64, customerId int64) error {
	updateItemQuery := `
		UPDATE cart_items
		SET quantity = quantity + 1,
			price = (
				SELECT price FROM products WHERE id = (
					SELECT product_id FROM cart_items WHERE id = ?
				)
			),
			updated_at = NOW(),
			updated_by = ?
		WHERE id = ?
	`

	if _, err := r.ExecContext(ctx, updateItemQuery, cartItemId, customerId, cartItemId); err != nil {
		return err
	}

	updateTotalQuery := `
		UPDATE carts
		SET total_amount = (
			SELECT COALESCE(SUM(ci.quantity * ci.price), 0)
			FROM cart_items ci
			WHERE ci.cart_id = (
				SELECT cart_id FROM cart_items WHERE id = ?
			)
			  AND ci.deleted_at IS NULL
		)
		WHERE id = (
			SELECT cart_id FROM cart_items WHERE id = ?
		)
	`

	if _, err := r.ExecContext(ctx, updateTotalQuery, cartItemId, cartItemId); err != nil {
		return err
	}

	return nil
}

func (r Repository) DeleteCartItem(ctx context.Context, cartItemId int64, customerId int64) error {
	deleteItemQuery := `
		UPDATE cart_items
		SET deleted_at = NOW(),
			deleted_by = ?
		WHERE id = ?
	`

	if _, err := r.ExecContext(ctx, deleteItemQuery, customerId, cartItemId); err != nil {
		return err
	}

	updateTotalQuery := `
		UPDATE carts
		SET total_amount = (
			SELECT COALESCE(SUM(ci.quantity * ci.price), 0)
			FROM cart_items ci
			WHERE ci.cart_id = (
				SELECT cart_id FROM cart_items WHERE id = ?
			)
			  AND ci.deleted_at IS NULL
		)
		WHERE id = (
			SELECT cart_id FROM cart_items WHERE id = ?
		)
	`

	if _, err := r.ExecContext(ctx, updateTotalQuery, cartItemId, cartItemId); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetList(ctx context.Context, customerId int64) ([]cart.Get, error) {
	query := `
		SELECT
			c.id AS cart_id,
			c.customer_id,
			c.total_amount,
			ci.id AS cart_item_id,
			ci.quantity AS item_quantity,
			ci.price AS item_price,
			ci.created_at AS item_created_at,
			p.id AS product_id,
			p.name,
			p.description,
			p.images,
			p.views_count,
			p.rating
		FROM carts c
		LEFT JOIN cart_items ci ON c.id = ci.cart_id AND ci.deleted_at IS NULL
		LEFT JOIN products p ON ci.product_id = p.id AND p.deleted_at IS NULL AND p.status = true
		WHERE c.customer_id = ? AND c.deleted_at IS NULL
		ORDER BY c.id, ci.id desc
	`

	rows, err := r.QueryContext(ctx, query, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cartMap := make(map[int64]*cart.Get)

	for rows.Next() {
		var (
			cartId      int64
			customerID  int64
			totalAmount float64

			itemID       sql.NullInt64
			itemQuantity sql.NullInt64
			itemPrice    sql.NullFloat64

			productID   sql.NullInt64
			name        sql.NullString
			description sql.NullString
			imagesJSON  []byte
			viewsCount  sql.NullInt64
			rating      sql.NullInt64
		)

		err := rows.Scan(
			&cartId,
			&customerID,
			&totalAmount,
			&itemID,
			&itemQuantity,
			&itemPrice,
			new(any),
			&productID,
			&name,
			&description,
			&imagesJSON,
			&viewsCount,
			&rating,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := cartMap[cartId]; !ok {
			cartMap[cartId] = &cart.Get{
				Id:          cartId,
				CustomerId:  customerID,
				TotalAmount: totalAmount,
				Items:       []cart.Item{},
			}
		}

		if itemID.Valid && productID.Valid {
			var imagesArray []entity.File
			if len(imagesJSON) > 0 {
				if err := json.Unmarshal(imagesJSON, &imagesArray); err != nil {
					log.Println("Failed to unmarshal images:", err)
				}
			}

			cartMap[cartId].Items = append(cartMap[cartId].Items, cart.Item{
				Id:          itemID.Int64,
				Name:        name.String,
				Description: description.String,
				Images:      &imagesArray,
				Quantity:    int(itemQuantity.Int64),
				Price:       itemPrice.Float64,
				Rating:      int8(rating.Int64),
				ViewsCount:  int(viewsCount.Int64),
				ProductId:   productID.Int64,
			})
		}
	}

	var result []cart.Get
	for _, c := range cartMap {
		result = append(result, *c)
	}

	return result, nil
}

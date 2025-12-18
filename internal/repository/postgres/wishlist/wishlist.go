package wishlist

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"main/internal/entity"
	"main/internal/services/wishlist"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) GetList(ctx context.Context, userId int64) ([]wishlist.GetList, error) {
	log.Println("userId", userId)
	query := `
        SELECT
            wl.id AS wishlist_id,
            wl.customer_id,
            wl.created_at AS wishlist_created_at,
            wli.id AS item_id,
            wli.product_id,
            wli.created_at AS item_created_at,
            p.name,
            p.price,
            p.images,
            p.views_count,
            p.discount_percent,
            p.category_id,
            p.rating,
            p.description
        FROM wishlists wl
        LEFT JOIN wishlist_items wli ON wl.id = wli.wishlist_id
        LEFT JOIN products p ON wli.product_id = p.id
        WHERE wl.customer_id = ? AND wl.deleted_at IS NULL AND wl.status = true AND (p.deleted_at IS NULL OR p.id IS NULL AND p.status = true) AND wli.deleted_at IS NULL
    `

	rows, err := r.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wishlistMap := make(map[int64]*wishlist.GetList)

	for rows.Next() {
		var wlID int64
		var customerID int64
		var wlCreatedAt time.Time
		var item wishlist.Item
		var itemID sql.NullInt64
		var productID sql.NullInt64
		var itemCreatedAt sql.NullTime
		var name, description sql.NullString
		var price sql.NullFloat64
		var imagesJSON []byte
		var viewsCount, discountPercent, categoryId sql.NullInt64
		var rating sql.NullInt64

		err := rows.Scan(
			&wlID, &customerID, &wlCreatedAt,
			&itemID, &productID, &itemCreatedAt,
			&name, &price, &imagesJSON, &viewsCount, &discountPercent,
			&categoryId, &rating, &description,
		)
		if err != nil {
			return nil, err
		}

		if _, ok := wishlistMap[wlID]; !ok {
			wishlistMap[wlID] = &wishlist.GetList{
				Id:        wlID,
				UserId:    customerID,
				CreatedAt: wlCreatedAt,
				Items:     []wishlist.Item{},
			}
		}

		var imagesArray []entity.File
		if len(imagesJSON) > 0 {
			if err := json.Unmarshal(imagesJSON, &imagesArray); err != nil {
				log.Println("Failed to unmarshal images:", err)
			}
		}

		if itemID.Valid {
			item = wishlist.Item{
				Id:              itemID.Int64,
				ProductId:       productID.Int64,
				Name:            name.String,
				Price:           price.Float64,
				Description:     description.String,
				ViewsCount:      viewsCount.Int64,
				DiscountPercent: int64(discountPercent.Int64),
				CategoryId:      categoryId.Int64,
				Rating:          int8(rating.Int64),
				CreatedAt:       itemCreatedAt.Time,
				Images:          &imagesArray,
			}
			wishlistMap[wlID].Items = append(wishlistMap[wlID].Items, item)
		}
	}

	var result []wishlist.GetList
	for _, wl := range wishlistMap {
		result = append(result, *wl)
	}

	return result, nil
}

func (r *Repository) Delete(ctx context.Context, wishlistItemId, customerId int64) error {
	log.Println("wishlistItemId", wishlistItemId, "customerId", customerId)
	query := `
		UPDATE wishlist_items wi
		SET deleted_at = NOW(), deleted_by = ?
		FROM wishlists wl
		WHERE wi.id = ?
		  AND wi.wishlist_id = wl.id
		  AND wl.customer_id = ?
		  AND wi.deleted_at IS NULL
	`

	if _, err := r.ExecContext(ctx, query, customerId, wishlistItemId, customerId); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, wishlist wishlist.Create, customerId int64) (int64, error) {
	var wishlistId int64

	getQuery := `
		SELECT id FROM wishlists
		WHERE customer_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.QueryRowContext(ctx, getQuery, customerId).Scan(&wishlistId)

	if err != nil {
		if err == sql.ErrNoRows {

			createQuery := `
				INSERT INTO wishlists (
					customer_id,
					created_at,
					created_by
				)
				VALUES (?, NOW(), ?)
				RETURNING id
			`

			if err := r.QueryRowContext(
				ctx,
				createQuery,
				customerId,
				customerId,
			).Scan(&wishlistId); err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	itemQuery := `
		INSERT INTO wishlist_items (
			wishlist_id,
			product_id,
			created_at,
			created_by
		)
		VALUES (?, ?, NOW(), ?)
		RETURNING id
	`

	var itemId int64
	if err := r.QueryRowContext(
		ctx,
		itemQuery,
		wishlistId,
		wishlist.ProductId,
		customerId,
	).Scan(&itemId); err != nil {
		return 0, err
	}

	return wishlistId, nil
}

package product

import (
	"context"
	"encoding/json"
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

func (r Repository) Create(ctx context.Context, data product.Create, adminId int64) (int64, error) {
	var newProductid int64
	var productParamValueId int64
	query := `INSERT INTO products (name, description, price, stock_quantity, category_id, discount_percent, images, created_by, status, created_at, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	if err := r.QueryRowContext(ctx, query, data.Name, data.Description, data.Price, data.StockQuantity, data.CategoryId, data.DiscountPercent, data.Images, adminId, false, time.Now(), adminId).Scan(&newProductid); err != nil {
		return 0, err
	}

	query = `INSERT INTO product_param_values (product_id, param_id, value_id, created_by, created_at) VALUES (?, ?, ?, ?, NOW()) RETURNING id`

	for _, param := range data.ParamSelected {
		for _, value := range param.ValueIDs {
			if err := r.QueryRowContext(ctx, query, newProductid, param.ParamID, value, adminId).Scan(&productParamValueId); err != nil {
				return 0, err
			}
		}
	}
	return newProductid, nil
}

func (r Repository) GetById(ctx context.Context, id int64, lang string) (product.GetById, error) {
	var data product.GetById

	query := fmt.Sprintf(`
			SELECT
			id,
			name,
			description,
			price,
			stock_quantity,
			rating_avg,
			seller_id,
			category_id,
			views_count,
			discount_percent,
			images,
			created_at
			FROM products
			WHERE id = %d AND deleted_at IS NULL
		`, id)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return product.GetById{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &data); err != nil {
		return product.GetById{}, err
	}

	params, err := r.getProductParams(ctx, data.Id)
	if err != nil {
		return product.GetById{}, err
	}

	data.Params = params

	return data, nil
}

func (r Repository) getProductParams(ctx context.Context, productID int64) ([]product.ParamWithValues, error) {
	query := `
		SELECT param_id, value_id
		FROM product_param_values
		WHERE product_id = ?
	`

	rows, err := r.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	paramMap := make(map[int64][]int64)

	for rows.Next() {
		var paramID, valueID int64
		if err := rows.Scan(&paramID, &valueID); err != nil {
			return nil, err
		}
		paramMap[paramID] = append(paramMap[paramID], valueID)
	}

	var result []product.ParamWithValues
	for pid, values := range paramMap {
		result = append(result, product.ParamWithValues{
			ParamID:  pid,
			ValueIDs: values,
		})
	}

	return result, nil
}

func (r Repository) GetList(ctx context.Context, filter entity.Filter) ([]product.Get, int, error) {
	var data []product.Get
	var limitQuery, offsetQuery string

	whereQuery := "WHERE p.deleted_at IS NULL"
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
		parts := strings.Split(*filter.Order, "+")
		log.Println("filter.Order repo", parts)

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
		searchQuery := fmt.Sprintf(" AND (p.name ->> '%s' ILIKE '%%%s%%' OR c.name ->> '%s' ILIKE '%%%s%%')", *filter.Language, *filter.Search, *filter.Language, *filter.Search)
		whereQuery += searchQuery
	}

	query := fmt.Sprintf(`
		SELECT
			p.id,
			p.name ->> '%s' as name,
			p.description ->> '%s' as description,
			p.price,
			p.stock_quantity,
			p.rating_avg,
			p.seller_id,
			p.category_id,
			p.views_count,
			p.discount_percent,
			p.status,
			p.images,
			c.name ->> '%s' as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id AND c.deleted_at IS NULL
		%s
		%s
		%s
		%s
	`, *filter.Language, *filter.Language, *filter.Language, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &data); err != nil {
		return nil, 0, err
	}

	countQuery := `SELECT COUNT(p.id) FROM products p WHERE p.deleted_at IS NULL`
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

func (r Repository) UpdateProduct(ctx context.Context, productId int, data product.Update, userId int64) error {
	var product entity.Product

	var nameJSON []byte
	var descriptionJSON []byte
	var imagesJSON []byte

	query := `SELECT id, name, description, stock_quantity, status, seller_id, category_id, discount_percent, images , price FROM products WHERE id = ?  AND deleted_at IS NULL`
	if err := r.QueryRowContext(ctx, query, productId).Scan(&product.Id, &nameJSON, &descriptionJSON, &product.StockQuantity, &product.Status, &product.SellerId, &product.CategoryId, &product.DiscountPercent, &imagesJSON, &product.Price); err != nil {
		return err
	}

	if err := json.Unmarshal(nameJSON, &product.Name); err != nil {
		return err
	}

	if err := json.Unmarshal(descriptionJSON, &product.Description); err != nil {
		return err
	}

	if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
		return err
	}

	if data.Name != nil {
		if data.Name.Uz != nil {
			product.Name.Uz = data.Name.Uz
		}

		if data.Name.Ru != nil {
			product.Name.Ru = data.Name.Ru
		}

		if data.Name.En != nil {
			product.Name.En = data.Name.En
		}
	}

	if data.Description != nil {
		if data.Description.Uz != nil {
			product.Description.Uz = data.Description.Uz
		}

		if data.Description.Ru != nil {
			product.Description.Ru = data.Description.Ru
		}

		if data.Description.En != nil {
			product.Description.En = data.Description.En
		}
	}

	if data.Status != nil {
		product.Status = *data.Status
	}

	if data.Price != nil {
		product.Price = *data.Price
	}

	if data.DiscountPercent != nil {
		product.DiscountPercent = *data.DiscountPercent
	}

	if data.Images != nil {
		product.Images = *data.Images
	}

	if data.SellerId != nil {
		product.SellerId = *data.SellerId
	}

	if data.CategoryId != nil {
		product.CategoryId = *data.CategoryId
	}

	if data.StockQuantity != nil {
		product.StockQuantity = *data.StockQuantity
	}

	if data.Images != nil {
		product.Images = *data.Images
	}

	var imagesJson string
	if product.Images != nil {
		b, err := json.Marshal(product.Images)
		if err != nil {
			return err
		}
		imagesJson = string(b)
	} else {
		imagesJson = "[]"
	}

	query = `UPDATE products SET name = ?, description = ?, price = ?, discount_percent = ?, images = ?, seller_id = ?, category_id = ?, status = ?, updated_at = NOW(), updated_by = ?, stock_quantity = ? WHERE id = ? AND deleted_at IS NULL`
	if _, err := r.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.DiscountPercent, imagesJson, product.SellerId, product.CategoryId, product.Status, userId, product.StockQuantity, productId); err != nil {
		return err
	}

	return nil
}

func (r Repository) Delete(ctx context.Context, productId, userId int64) error {
	query := `UPDATE products SET deleted_at = NOW(), deleted_by = ? WHERE id = ? AND deleted_at IS NULL `
	if _, err := r.ExecContext(ctx, query, userId, productId); err != nil {
		return err
	}

	return nil
}

func (r Repository) GetByIdDetail(ctx context.Context, id int64) (product.GetById, error) {
	var data product.GetById

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
			WHERE p.id = ? AND p.deleted_at IS NULL
		`
	rows, err := r.QueryContext(ctx, query, id)
	if err != nil {
		return product.GetById{}, err
	}
	defer rows.Close()

	if err := r.ScanRows(ctx, rows, &data); err != nil {
		return product.GetById{}, err
	}

	return data, nil
}

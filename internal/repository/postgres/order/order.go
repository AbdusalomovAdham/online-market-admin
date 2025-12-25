package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"main/internal/entity"
	order "main/internal/services/order"

	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Repository struct {
	*bun.DB
}

func NewRepository(DB *bun.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, order order.Create, customerId int64) error {
	var totalAmount float64
	var orderId int64

	for i := range order.Items {
		item := &order.Items[i]
		totalAmount += item.Price * float64(item.Quantity)
	}
	order.TotalAmount = totalAmount

	if order.DeliveryDate == "" {
		order.DeliveryDate = time.Now().Add(72 * time.Hour).Format("2006-01-02")
	}

	query := `
			INSERT INTO orders (
					order_status,
					payment_status,
					delivery_date,
					total_amount,
					customer_id,
					created_at,
					created_by
			) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := r.QueryRowContext(
		ctx,
		query,
		order.OrderStatus,
		order.PaymentId,
		order.DeliveryDate,
		order.TotalAmount,
		customerId,
		time.Now(),
		customerId,
	).Scan(&orderId)
	if err != nil {
		return err
	}

	itemQuery := `
			INSERT INTO order_items (
					order_id,
					product_id,
					quantity,
					price,
					total,
					created_at,
					created_by
			) VALUES (?, ?, ?, ?, ?, ?, ?)`

	for _, item := range order.Items {
		total := item.Price * float64(item.Quantity)
		_, err := r.ExecContext(
			ctx,
			itemQuery,
			orderId,
			item.ProductId,
			item.Quantity,
			item.Price,
			total,
			time.Now(),
			customerId,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetList(ctx context.Context, userId int64, filter entity.Filter) ([]order.Get, int64, error) {
	var limitQuery, offsetQuery string

	whereQuery := fmt.Sprintf("WHERE o.deleted_at IS NULL AND o.customer_id = %d", userId)

	if filter.Limit != nil {
		limitQuery = fmt.Sprintf("LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery = fmt.Sprintf("OFFSET %d", *filter.Offset)
	}

	orderQuery := "ORDER BY o.id DESC"
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
		  	SELECT o.id, os.name ->> '%s' AS order_status, ps.name ->> '%s' AS payment_status, o.order_status_id, o.payment_id, o.delivery_date, o.total_amount
		    FROM orders o
		 	LEFT JOIN order_statuses os ON os.id = o.order_status_id
			LEFT JOIN payments ps ON ps.id = o.payment_id
		    %s %s %s %s
	`, *filter.Language, *filter.Language, whereQuery, orderQuery, limitQuery, offsetQuery)

	var orders []order.Get

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var o order.Get
		if err := rows.Scan(&o.Id, &o.OrderStatus, &o.PaymentStatus, &o.OrderStatusId, &o.PaymentId, &o.DeliveryDate, &o.TotalAmount); err != nil {
			return nil, 0, err
		}
		o.Items = []order.GetItems{}
		orders = append(orders, o)
	}

	if len(orders) == 0 {
		return nil, 0, errors.New("No orders found")
	}

	ordersId := make([]int64, len(orders))
	orderMap := make(map[int64]*order.Get)
	for i := range orders {
		ordersId[i] = orders[i].Id
		orderMap[orders[i].Id] = &orders[i]
	}

	itemsQuery := `
				SELECT
			    oi.id,
			    oi.order_id,
			   	p.name ->> ? AS name,
				p.description ->> ? AS description,
			    COALESCE(p.price, 0) AS price,
				COALESCE(p.images, '[]') AS images,
			    oi.quantity,
			    COALESCE(p.rating_avg, 0) AS rating
			FROM order_items oi
			LEFT JOIN products p ON p.id = oi.product_id
			WHERE oi.order_id = ANY(?) AND p.deleted_at IS NULL AND p.status = true AND oi.deleted_at IS NULL
    `

	itemRows, err := r.QueryContext(ctx, itemsQuery, filter.Language, filter.Language, pq.Array(ordersId))
	if err != nil {
		return nil, 0, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item struct {
			Id          int64
			OrderId     int64
			Name        string
			Description string
			Price       float64
			Quantity    int
			Rating      float32
			Images      []byte
		}

		if err := itemRows.Scan(&item.Id, &item.OrderId, &item.Name, &item.Description, &item.Price, &item.Images, &item.Quantity, &item.Rating); err != nil {
			return nil, 0, err
		}

		var images []entity.File
		if len(item.Images) > 0 {
			if err := json.Unmarshal(item.Images, &images); err != nil {
				log.Println("Failed to unmarshal images:", err)
			}
		}

		if o, ok := orderMap[item.OrderId]; ok {
			o.Items = append(o.Items, order.GetItems{
				Id:          item.Id,
				Name:        item.Name,
				Description: item.Description,
				Price:       item.Price,
				Quantity:    item.Quantity,
				Rating:      item.Rating,
				Images:      &images,
			})
		}
	}

	countQuery := `SELECT COUNT(o.id) FROM orders o WHERE o.deleted_at IS NULL`

	countRows, err := r.QueryContext(ctx, countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()

	count := int64(0)
	if err = r.ScanRows(ctx, countRows, &count); err != nil {
		return nil, 0, fmt.Errorf("select category count: %w", err)
	}

	return orders, count, nil
}

func (r Repository) GetById(ctx context.Context, orderId, userId int64) (order.Get, error) {
	var o order.Get

	query := `
        SELECT id, order_status, payment_status, delivery_date, total_amount
        FROM orders
        WHERE id = ? AND customer_id = ? AND deleted_at IS NULL
    `
	if err := r.QueryRowContext(ctx, query, orderId, userId).Scan(
		&o.Id,
		&o.OrderStatus,
		&o.PaymentStatus,
		&o.DeliveryDate,
		&o.TotalAmount,
	); err != nil {
		return o, err
	}

	o.Items = []order.GetItems{}

	itemsQuery := `
        SELECT
            oi.id,
            oi.order_id,
            COALESCE(p.name, '') AS name,
            COALESCE(p.description, '') AS description,
            COALESCE(p.price, 0) AS price,
            COALESCE(p.images, '[]') AS images,
            oi.quantity,
            COALESCE(p.rating, 0) AS rating
        FROM order_items oi
        LEFT JOIN products p ON p.id = oi.product_id
        WHERE oi.order_id = ? AND p.deleted_at IS NULL AND p.status = true AND oi.deleted_at IS NULL
    `
	itemRows, err := r.QueryContext(ctx, itemsQuery, orderId)
	if err != nil {
		return o, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item struct {
			Id          int64
			OrderId     int64
			Name        string
			Description string
			Price       float64
			Quantity    int
			Rating      float32
			Images      []byte
		}

		if err := itemRows.Scan(
			&item.Id,
			&item.OrderId,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.Images,
			&item.Quantity,
			&item.Rating,
		); err != nil {
			return o, err
		}

		var images []entity.File
		if len(item.Images) > 0 {
			if err := json.Unmarshal(item.Images, &images); err != nil {
				log.Println("Failed to unmarshal images:", err)
			}
		}

		o.Items = append(o.Items, order.GetItems{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Rating:      item.Rating,
			Images:      &images,
		})
	}

	return o, nil
}

func (r Repository) Delete(ctx context.Context, orderId int64, userId int64) error {
	query := `
        UPDATE orders
        SET deleted_at = NOW(), deleted_by = ?
        WHERE id = ?
    `
	_, err := r.ExecContext(ctx, query, userId, orderId)
	if err != nil {
		return err
	}

	return nil
}

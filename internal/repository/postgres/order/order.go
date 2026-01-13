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

func (r *Repository) GetList(ctx context.Context, userId int64, filter entity.Filter) ([]order.GetList, int64, error) {
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
		SELECT
			o.id,
			os.id,
			p.id,
			os.key AS order_status_key,
			os.name ->> '%s' AS order_status_value,
			p.key AS payment_status_key,
			p.name ->> '%s' AS payment_status_value,
			o.delivery_date,
			o.total_amount,
			COUNT(oi.id) AS items_count,
			o.created_at,
			TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS customer_name
		FROM orders o
		LEFT JOIN payments p ON p.id = o.payment_id
		LEFT JOIN order_statuses os ON os.id = o.order_status_id
		LEFT JOIN order_items oi ON oi.order_id = o.id
		LEFT JOIN users u ON u.id = o.customer_id
		%s
		GROUP BY
			o.id,
			os.id,
			p.id,
			os.key,
			os.name,
			p.key,
			p.name,
			os.name,
			p.name,
			o.delivery_date,
			o.total_amount,
			u.first_name,
			u.last_name
		%s %s %s`,
		*filter.Language,
		*filter.Language,
		whereQuery,
		orderQuery,
		limitQuery,
		offsetQuery,
	)

	var orders []order.GetList
	var (
		orderStatusId    int
		paymentStatusId  int
		orderStatusKey   string
		orderStatusValue string
		paymentStatusKey *string
		paymentStatusVal *string
	)

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var o order.GetList
		if err := rows.Scan(
			&o.Id,
			&orderStatusId,
			&paymentStatusId,
			&orderStatusKey,
			&orderStatusValue,
			&paymentStatusKey,
			&paymentStatusVal,
			&o.DeliveryDate,
			&o.TotalAmount,
			&o.ItemsCount,
			&o.CreatedAt,
			&o.CustomerName,
		); err != nil {
			return nil, 0, err
		}

		o.OrderStatus = order.OrderStatus{
			Id:    int64(orderStatusId),
			Key:   orderStatusKey,
			Value: orderStatusValue,
		}

		if paymentStatusKey != nil {
			o.PaymentStatus = order.PaymentStatus{
				Id:    int64(paymentStatusId),
				Key:   *paymentStatusKey,
				Value: *paymentStatusVal,
			}
		} else {
			o.PaymentStatus = order.PaymentStatus{
				Key:   "pending",
				Value: "Not paid",
			}
		}

		orders = append(orders, o)
	}

	if len(orders) == 0 {
		return nil, 0, errors.New("No orders found")
	}

	ordersId := make([]int64, len(orders))
	orderMap := make(map[int64]*order.GetList)
	for i := range orders {
		ordersId[i] = orders[i].Id
		orderMap[orders[i].Id] = &orders[i]
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

func (r Repository) GetById(ctx context.Context, orderId int64, filter entity.Filter) (order.Get, error) {
	var o order.Get
	var orderStatusName, paymentStatusName, orderStatusKey, paymentStatusKey string
	var orderStautsId, paymentStatusId int64
	query := `
        SELECT
	        o.id,
	        o.delivery_date,
	        o.total_amount,
	        p.key,
	        os.key,
	        os.name ->> ? ,
	        p.name ->> ?,
			p.id,
			os.id,
			TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))) AS customer_name,
			u.email,
			u.phone_number,
			d.name ->> ?,
			r.name ->> ?
        FROM orders o
        LEFT JOIN order_statuses os ON os.id = o.order_status_id
        LEFT JOIN users u ON u.id = o.customer_id
        LEFT JOIN payments p ON p.id = o.payment_id
        LEFT JOIN districts d ON d.id = u.district_id
        LEFT JOIN regions r ON r.id = d.region_id
        WHERE o.id = ? AND o.deleted_at IS NULL
    `

	if err := r.QueryRowContext(ctx, query, filter.Language, filter.Language, filter.Language, filter.Language, orderId).Scan(
		&o.Id,
		&o.DeliveryDate,
		&o.TotalAmount,
		&paymentStatusKey,
		&orderStatusKey,
		&orderStatusName,
		&paymentStatusName,
		&orderStautsId,
		&paymentStatusId,
		&o.CustomerName,
		&o.Email,
		&o.PhoneNumber,
		&o.DistrictName,
		&o.RegionName,
	); err != nil {
		return o, err
	}

	o.OrderStatus.Id = orderStautsId
	o.PaymentStatus.Id = paymentStatusId
	o.OrderStatus.Value = orderStatusName
	o.PaymentStatus.Value = paymentStatusName
	o.OrderStatus.Key = orderStatusKey
	o.PaymentStatus.Key = paymentStatusKey

	o.Items = []order.GetItems{}

	itemsQuery := `
		SELECT
		    oi.id,
		    oi.order_id,
		    COALESCE(p.name ->> ?, ''),
		    COALESCE(p.description ->> ?, ''),
		    COALESCE(p.price, 0),
		    p.images,
		    oi.quantity,
		    COALESCE(p.rating_avg, 0),
		    p.discount_percent
		FROM order_items oi
		INNER JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = ? AND p.deleted_at IS NULL AND p.status = true AND oi.deleted_at IS NULL`

	itemRows, err := r.QueryContext(ctx, itemsQuery, filter.Language, filter.Language, orderId)
	if err != nil {
		return o, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item struct {
			Id              int64
			OrderId         int64
			Name            string
			Description     string
			Price           float64
			Quantity        int
			Rating          float32
			Images          []byte
			DiscountPercent int
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
			&item.DiscountPercent,
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
			Id:              item.Id,
			Name:            item.Name,
			Description:     item.Description,
			Price:           item.Price,
			Quantity:        item.Quantity,
			Rating:          item.Rating,
			Images:          &images,
			DiscountPercent: item.DiscountPercent,
			OrderId:         int64(item.OrderId),
		})
	}

	countQuery := `
    SELECT COUNT(oi.id)
    FROM order_items oi
    LEFT JOIN products p ON p.id = oi.product_id
    WHERE oi.order_id = ?
      AND p.deleted_at IS NULL
      AND p.status = true
      AND oi.deleted_at IS NULL
      `

	if err := r.QueryRowContext(ctx, countQuery, orderId).Scan(&o.ItemsCount); err != nil {
		return o, err
	}

	return o, nil
}

func (r Repository) Update(ctx context.Context, updateData order.Update, orderId int64, adminId int64) error {
	var order entity.Order
	query := `SELECT
				id,
				payment_id,
				order_status_id,
				delivery_date
				FROM orders
				WHERE id = ?
				`

	row := r.QueryRowContext(ctx, query, orderId)

	err := row.Scan(&order.Id, &order.PaymentId, &order.OrderStatusId, &order.DeliveryDate)
	if err != nil {
		return err
	}

	if updateData.OrderStatus != nil {
		if *updateData.OrderStatus > 8 && *updateData.OrderStatus < 1 {
			return errors.New("invalid order status")
		}
		order.OrderStatusId = *updateData.OrderStatus
	}

	if updateData.PaymentStatus != nil {
		if *updateData.PaymentStatus > 5 && *updateData.PaymentStatus < 1 {
			return errors.New("invalid payment status")
		}
		order.PaymentId = *updateData.PaymentStatus
	}

	if updateData.DeliveryDate != nil {
		order.DeliveryDate = *updateData.DeliveryDate
	}

	updateQuery := `
        UPDATE orders
        SET payment_id = ?, order_status_id = ?, delivery_date = ?, updated_at = NOW(), updated_by = ?
        WHERE id = ?
    `
	_, err = r.ExecContext(ctx, updateQuery, order.PaymentId, order.OrderStatusId, order.DeliveryDate, adminId, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) DeleteOrderItem(ctx context.Context, itemId, adminId int64) error {
	query := `
        UPDATE order_items
        SET deleted_at = NOW(), deleted_by = ?
        WHERE id = ?
    `
	_, err := r.ExecContext(ctx, query, adminId, itemId)
	if err != nil {
		return err
	}

	return nil
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

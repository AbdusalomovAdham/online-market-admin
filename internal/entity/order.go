package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders"`

	Id            int64       `json:"id" bun:"id,pk,autoincrement"`
	OrderStatus   string      `json:"order_status" bun:"order_status"`
	PaymentStatus string      `json:"payment_status" bun:"payment_status"`
	DeliveryDate  string      `json:"delivery_date" bun:"delivery_date"`
	TotalAmount   float64     `json:"total_amount" bun:"total_amount"`
	CustomerId    int64       `json:"customer_id" bun:"customer_id"`
	Items         []OrderItem `json:"items" bun:"items"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

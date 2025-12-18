package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items"`

	Id        int64   `json:"id" bun:"id,pk,autoincrement"`
	OrderId   int64   `json:"order_id" bun:"order_id"`
	ProductId int64   `json:"product_id" bun:"product_id"`
	Quantity  int64   `json:"quantity" bun:"quantity"`
	Price     float64 `json:"price" bun:"price"`
	Total     float64 `json:"total" bun:"total"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

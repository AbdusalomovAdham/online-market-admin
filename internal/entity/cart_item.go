package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type CartItem struct {
	bun.BaseModel `bun:"table:cart_items"`

	Id        int64   `json:"id" bun:"id,pk,autoincrement"`
	ProductId int64   `json:"product_id" bun:"product_id"`
	CartId    int64   `json:"cart_id" bun:"cart_id"`
	Status    bool    `json:"status" bun:"status"`
	Quantity  int     `json:"quantity" bun:"quantity"`
	Price     float64 `json:"price" bun:"price"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

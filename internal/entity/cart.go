package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Cart struct {
	bun.BaseModel `bun:"table:carts"`

	Id          int64   `json:"id" bun:"id,pk,autoincrement"`
	CustomerId  int64   `json:"customer_id" bun:"customer_id"`
	Status      bool    `json:"status" bun:"status"`
	TotalAmount float64 `json:"total_amount" bun:"total_amount"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

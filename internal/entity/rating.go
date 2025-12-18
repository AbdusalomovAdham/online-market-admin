package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Rating struct {
	bun.BaseModel `bun:"table:ratings"`

	Id         int64 `json:"id" bun:"id,pk,autoincrement"`
	CustomerId int64 `json:"customer_id" bun:"customer_id"`
	ProductId  int64 `json:"product_id" bun:"product_id"`
	Rating     int8  `json:"rating" bun:"rating" default:"0"`
	Status     bool  `json:"status" bun:"status" default:"true"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

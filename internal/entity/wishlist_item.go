package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type WishlistItem struct {
	bun.BaseModel `bun:"table:wishlist_items"`

	Id         int64 `json:"id" bun:"id,pk,autoincrement"`
	ProductId  int64 `json:"product_id" bun:"product_id"`
	WishlistId int64 `json:"wishlist_id" bun:"wishlist_id"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy int64      `json:"-" bun:"created_by"`
	UpdatedAt *time.Time `json:"-" bun:"updated_at" default:"null"`
	UpdatedBy *int64     `json:"-" bun:"updated_by" default:"null"`
	DeletedAt *time.Time `json:"-" bun:"deleted_at" default:"null"`
	DeletedBy *int64     `json:"-" bun:"deleted_by" default:"null"`
}

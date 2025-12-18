package wishlist

import (
	"main/internal/entity"
	"time"
)

type Create struct {
	ProductId int64 `json:"product_id" bun:"product_id"`
}

type CreateItem struct {
	ProductId  int64 `json:"product_id" bun:"product_id"`
	WishlistId int64 `json:"wishlist_id" bun:"wishlist_id"`
}

type GetList struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id" bun:"user_id"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	Items     []Item    `json:"items"`
}

type Item struct {
	Id              int64          `json:"id"`
	ProductId       int64          `json:"product_id" bun:"product_id"`
	Name            string         `json:"name" bun:"name"`
	Description     string         `json:"description" bun:"description"`
	Price           float64        `json:"price" bun:"price"`
	Images          *[]entity.File `json:"images" bun:"images"`
	ViewsCount      int64          `json:"views_count" bun:"views_count"`
	DiscountPercent int64          `json:"discount_percent" bun:"discount_percent"`
	CategoryId      int64          `json:"category_id" bun:"category_id"`
	Rating          int8           `json:"rating" bun:"rating"`
	CreatedAt       time.Time      `json:"created_at" bun:"created_at"`
}

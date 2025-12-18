package products

import (
	"main/internal/entity"
	"time"
)

type Create struct {
	Name            string         `json:"name" form:"name"`
	Description     *string        `json:"description" form:"description"`
	Price           int64          `json:"price" form:"price"`
	StockQuantity   *int64         `json:"stock_quantity" form:"stock_quantity"`
	CategoryId      int64          `json:"category_id" form:"category_id"`
	DiscountPercent *int8          `json:"discount_percent" form:"discount_percent"`
	Images          *[]entity.File `json:"images"`
}

type Get struct {
	Id              int64          `json:"id" bun:"id"`
	Name            *string        `json:"name" bun:"name"`
	Description     *string        `json:"description" form:"description" bun:"description"`
	Price           float64        `json:"price" form:"price" bun:"price"`
	StockQuantity   *int64         `json:"stock_quantity" form:"stock_quantity" bun:"stock_quantity"`
	CategoryId      int64          `json:"category_id" form:"category_id" bun:"category_id"`
	DiscountPercent *int8          `json:"discount_percent" form:"discount_percent" bun:"discount_percent"`
	RatingAvg       float32        `json:"rating_avg" bun:"rating_avg"`
	SellerId        int64          `json:"seller_id" bun:"seller_id"`
	ViewsCount      int64          `json:"views_count" bun:"views_count"`
	IsWishlist      bool           `json:"is_wishlist" bun:"is_wishlist"`
	FirstName       string         `json:"first_name" bun:"first_name"`
	LastName        string         `json:"last_name" bun:"last_name"`
	Images          *[]entity.File `json:"images" bun:"images,type:jsonb"`
	Avatar          string         `json:"avatar" bun:"avatar"`
	CreatedAt       time.Time      `json:"created_at" bun:"created_at"`
}

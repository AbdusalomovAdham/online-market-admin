package order

import (
	"main/internal/entity"
)

type Create struct {
	OrderStatus   string  `json:"order_status"`
	PaymentStatus string  `json:"payment_status"`
	DeliveryDate  string  `json:"delivery_date"`
	TotalAmount   float64 `json:"total_amount"`
	Items         []Item  `json:"items"`
}

type Item struct {
	ProductId int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Get struct {
	Id            int64      `json:"id"`
	OrderStatus   string     `json:"order_status"`
	PaymentStatus string     `json:"payment_status"`
	DeliveryDate  string     `json:"delivery_date"`
	TotalAmount   float64    `json:"total_amount"`
	Items         []GetItems `json:"items"`
}

type GetItems struct {
	Id          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      *[]entity.File `json:"images"`
	Quantity    int            `json:"quantity"`
	Rating      int8           `json:"rating"`
	Price       float64        `json:"price"`
}

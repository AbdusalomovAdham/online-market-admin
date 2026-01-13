package order

import (
	"main/internal/entity"
	"time"
)

type Create struct {
	OrderStatus  string  `json:"order_status"`
	PaymentId    string  `json:"payment_id"`
	DeliveryDate string  `json:"delivery_date"`
	TotalAmount  float64 `json:"total_amount"`
	Items        []Item  `json:"items"`
}

type Item struct {
	ProductId int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type GetList struct {
	Id            int64         `json:"id"`
	OrderStatus   OrderStatus   `json:"order_status"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	DeliveryDate  string        `json:"delivery_date"`
	TotalAmount   float64       `json:"total_amount"`
	CustomerName  string        `json:"customer_name"`
	ItemsCount    int           `json:"items_count"`
	CreatedAt     time.Time     `json:"created_at"`
}

type Get struct {
	Id            int64         `json:"id"`
	OrderStatus   OrderStatus   `json:"order_status" bun:"-"`
	PaymentStatus PaymentStatus `json:"payment_status" bun:"-"`
	DeliveryDate  string        `json:"delivery_date"`
	TotalAmount   float64       `json:"total_amount"`
	Items         []GetItems    `json:"items"`
	CustomerName  string        `json:"customer_name"`
	ItemsCount    int           `json:"items_count"`
	CreatedAt     time.Time     `json:"created_at"`
	Email         string        `json:"email"`
	PhoneNumber   string        `json:"phone_number"`
	DistrictName  string        `json:"district_name"`
	RegionName    string        `json:"region_name"`
}

type PaymentStatus struct {
	Id    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OrderStatus struct {
	Id    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetItems struct {
	Id              int64          `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	Images          *[]entity.File `json:"images"`
	Quantity        int            `json:"quantity"`
	Rating          float32        `json:"rating"`
	Price           float64        `json:"price"`
	DiscountPercent int            `json:"discount_percent"`
	OrderId         int64          `json:"order_id"`
}

type Update struct {
	PaymentStatus *int    `json:"payment_status"`
	OrderStatus   *int    `json:"order_status"`
	DeliveryDate  *string `json:"delivery_date"`
}

package cart

import "main/internal/entity"

type Create struct {
	ProductId int64 `json:"product_id"`
}

type Get struct {
	Id          int64   `json:"id"`
	TotalAmount float64 `json:"total_amount"`
	CustomerId  int64   `json:"customer_id"`
	Items       []Item  `json:"items"`
}

type Item struct {
	Id          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ProductId   int64          `json:"product_id"`
	Images      *[]entity.File `json:"images"`
	Quantity    int            `json:"quantity"`
	Rating      int8           `json:"rating"`
	Price       float64        `json:"price"`
	ViewsCount  int            `json:"views_count"`
}

package order_status

import "main/internal/entity"

type Create struct {
	Name   *entity.Name `json:"name"`
	Status *bool        `json:"status" default:"true"`
}

type Get struct {
	Id        int64  `json:"id"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
}

type OrderStatusById struct {
	Id        int64        `json:"id"`
	Status    bool         `json:"status"`
	CreatedAt string       `json:"created_at"`
	Name      *entity.Name `json:"name"`
}

type Update struct {
	Name   *entity.Name `json:"name"`
	Status *bool        `json:"status"`
}

package param

import (
	"main/internal/entity"
	"time"
)

type Create struct {
	Name       entity.Name `json:"name" bun:"name"`
	CategoryId []int64     `json:"category_id" bun:"category_id"`
	Type       string      `json:"type" bun:"type"`
	Status     *bool       `json:"status" bun:"status"`
}

type ParamValue struct {
	Id        int64     `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name" bun:"name"`
	Status    bool      `json:"status" bun:"status"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
}

type ParamById struct {
	Id         int64       `json:"id" bun:"id"`
	Name       entity.Name `json:"name" bun:"name"`
	Status     bool        `json:"status" bun:"status"`
	CategoryId []int64     `json:"category_id" bun:"category_id"`
	Type       string      `json:"type" bun:"type"`
	CreatedAt  time.Time   `json:"created_at" bun:"created_at"`
}

type Get struct {
	Id           int64     `json:"id" bun:"id"`
	ParamName    string    `json:"param_name" bun:"param_name"`
	CategoryName string    `json:"category_name" bun:"category_name"`
	Status       bool      `json:"status" bun:"status"`
	Type         string    `json:"type" bun:"type"`
	CreatedAt    time.Time `json:"created_at" bun:"created_at"`
}

type GetByCategoryId struct {
	Id        int64     `json:"id" bun:"id"`
	ParamName string    `json:"param_name" bun:"param_name"`
	Status    bool      `json:"status" bun:"status"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
}

type UpdateParam struct {
	Name       *entity.Name `json:"name" bun:"name"`
	Type       *string      `json:"type" bun:"type"`
	CategoryId *[]int64     `json:"category_id" bun:"category_id"`
	Status     *bool        `json:"status" bun:"status"`
}

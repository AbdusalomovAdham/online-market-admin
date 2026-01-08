package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type CategoryParam struct {
	bun.BaseModel `bun:"table:category_params"`

	Id         int64   `json:"id" bun:"id,pk,autoincrement"`
	CategoryId []int64 `json:"category_id" bun:"category_id"`
	ParamId    int64   `json:"param_id" bun:"param_id"`
	Status     bool    `json:"status" bun:"status"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

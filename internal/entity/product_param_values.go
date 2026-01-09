package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type ProductParamValues struct {
	bun.BaseModel `bun:"table:product_param_values"`

	Id        int64 `json:"id"`
	ProductId int64 `json:"product_id"`
	ParamId   int64 `json:"param_id"`
	ValueId   int64 `json:"value_id"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

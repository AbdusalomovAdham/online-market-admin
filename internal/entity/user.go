package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id          int64   `json:"id" bun:"id,pk,autoincrement"`
	Avatar      *string `json:"avatar" bun:"avatar"`
	FirstName   string  `json:"first_name" bun:"first_name"`
	LastName    string  `json:"last_name" bun:"last_name"`
	PhoneNumber string  `json:"phone_number" bun:"phone_number"`
	Birthdate   *string `json:"birth_date" bun:"birth_date"`
	Email       *string `json:"email" bun:"email"`
	Status      *bool   `json:"status" bun:"status,default:false"`
	Role        int     `json:"role" bun:"role,default:2"`
	Password    *string `json:"password" bun:"password"`
	RegionId    *int64  `json:"region_id" bun:"region_id,default:null"`
	DistrictId  *int64  `json:"district_id" bun:"district_id,default:null"`

	CreatedBy int       `json:"created_by" bun:"created_by,default:null"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedBy int       `json:"updated_by" bun:"updated_by,default:null"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,default:null"`
	DeletedAt time.Time `json:"deleted_at" bun:"deleted_at,default:null"`
	DeletedBy int64     `json:"deleted_by" bun:"deleted_by,default:null"`
}

//Role

//1-admin
//2-seller
//3-customer
//4-distributor
//5-operator

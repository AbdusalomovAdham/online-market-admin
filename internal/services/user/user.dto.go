package user

import "time"

type Create struct {
	Avatar      *string `json:"avatar"`
	FirstName   string  `json:"first_name" form:"first_name"`
	LastName    string  `json:"last_name" form:"last_name"`
	PhoneNumber string  `json:"phone_number" form:"phone_number"`
	Password    *string `json:"password" form:"password"`
	Login       *string `json:"login" form:"login"`
	BirthDate   *string `json:"birth_date" form:"birth_date"`
	Email       *string `json:"email" form:"email"`
	Role        int     `json:"role" form:"role"`
	RegionID    *int    `json:"region_id" form:"region_id"`
	DistrictID  *int    `json:"district_id" form:"district_id"`
}

type Get struct {
	ID          int64     `json:"id"`
	Avatar      *string   `json:"avatar"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Login       *string   `json:"login"`
	BirthDate   *string   `json:"birth_date"`
	Email       *string   `json:"email"`
	Role        string    `json:"role"`
	RegionID    *int      `json:"region_id"`
	DistrictID  *int      `json:"district_id"`
	CreatedAt   time.Time `json:"created_at"`
	Password    *string   `json:"password" default:""`
}

type Update struct {
	Avatar      *string `json:"avatar"`
	FirstName   *string `json:"first_name" form:"first_name"`
	LastName    *string `json:"last_name" form:"last_name"`
	PhoneNumber *string `json:"phone_number" form:"phone_number"`
	Password    *string `json:"password" form:"password"`
	Login       *string `json:"login" form:"login"`
	BirthDate   *string `json:"birth_date" form:"birth_date"`
	Email       *string `json:"email" form:"email"`
	Role        *int    `json:"role" form:"role"`
	RegionId    *int    `json:"region_id" form:"region_id"`
	DistrictId  *int    `json:"district_id" form:"district_id"`
	Status      *bool   `json:"status" form:"status"`
}

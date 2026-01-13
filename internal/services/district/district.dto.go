package district

import "time"

type GetList struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	RegionId  int       `json:"region_id" bun:"region_id"`
}

type GetById struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

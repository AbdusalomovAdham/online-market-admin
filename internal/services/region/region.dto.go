package region

import "time"

type GetList struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
}

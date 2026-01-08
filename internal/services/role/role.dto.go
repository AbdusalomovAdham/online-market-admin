package role

import "time"

type Get struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `jons:"created_at"`
}

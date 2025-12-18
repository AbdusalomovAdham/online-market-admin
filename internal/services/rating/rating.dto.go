package rating

type Create struct {
	ProductId int64 `json:"product_id"`
	Rating    int   `json:"rating" bun:"rating" default:"1"`
}

package ws

import "context"

type Repository interface {
	GetDriverListByDistrictId(ctx context.Context, fromDistrictId, toDistrictId int) ([]int, error)
}

package param_value

import (
	"main/internal/entity"
	"time"
)

type Create struct {
	Name    entity.Name `json:"name"`
	ParamId int64       `json:"param_id"`
	Status  *bool       `json:"status"`
}

type Get struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	ParamId   int64     `json:"param_id"`
	ParamName string    `json:"param_name"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ParamValueById struct {
	Id        int64       `json:"id"`
	Name      entity.Name `json:"name"`
	ParamId   int64       `json:"param_id"`
	Status    bool        `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

type ParamValueByParamId struct {
	Id        int64       `json:"id"`
	Name      entity.Name `json:"name"`
	Status    bool        `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

type Update struct {
	Name    *entity.Name `json:"name"`
	Status  *bool        `json:"status"`
	ParamId *int64       `json:"param_id"`
}

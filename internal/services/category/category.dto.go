package category

import "main/internal/entity"

type Create struct {
	Name     *entity.Name `json:"name"`
	ParentId *int64       `json:"parent_id"`
	Status   *bool        `json:"status" default:"true"`
}

type Get struct {
	Id         int64  `json:"id"`
	Status     bool   `json:"status"`
	CreatedAt  string `json:"created_at"`
	ParentName string `json:"parent_name"`
	Name       string `json:"name"`
}

type ParamInfo struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type CategoryById struct {
	Id        int64        `json:"id" bson:"id"`
	Status    bool         `json:"status" bson:"status"`
	CreatedAt string       `json:"created_at" bson:"created_at"`
	ParentId  *int64       `json:"parent_id" bson:"parent_id"`
	Name      *entity.Name `json:"name" bson:"name"`
	Params    []ParamInfo  `json:"params" bson:"params"`
}

type Update struct {
	Name     *entity.Name `json:"name"`
	Status   *bool        `json:"status"`
	ParentId **int64      `json:"parent_id"`
}

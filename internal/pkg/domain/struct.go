package domain

import (
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
)

// Model
type Model struct {
	Id        uint   `json:"id"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

// ReqId
type ReqId struct {
	Id uint `json:"id" param:"id"`
}

// PaginateReq
type PaginateReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Field    string `json:"field"`
	Order    string `json:"order"`
}

func (r *PaginateReq) ConvertParams() {
	r.Field = stringUtils.SnakeCase(r.Field)
	r.Order = serverConfig.SortMap[r.Order]
}

// Response
type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type PageData struct {
	Result interface{} `json:"result"`

	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (d *PageData) Populate(result interface{}, total int64, page, pageSize int) {
	d.Result = result
	d.Total = int(total)
	d.Page = page
	d.PageSize = pageSize
}

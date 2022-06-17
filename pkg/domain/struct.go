package domain

import (
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
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

	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (d *PageData) Populate(result interface{}, total int64, page, pageSize int) {
	d.Result = result

	d.Pagination.Total = int(total)
	d.Pagination.Page = page
	d.Pagination.PageSize = pageSize
}

type NestedItem struct {
	Id       int           `json:"id"`
	Name     string        `json:"name"`
	Parent   int           `json:"parent"`
	Children []*NestedItem `json:"children"`
}

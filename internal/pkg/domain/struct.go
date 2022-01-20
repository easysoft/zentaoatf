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

// ErrMsg
type ErrMsg struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

var (
	NoErr         = ErrMsg{0, "请求成功"}
	CommonErr     = ErrMsg{1001, "请求失败"}
	NeedInitErr   = ErrMsg{2001, "前往初始化数据库"}
	AuthErr       = ErrMsg{4001, "会话超时，请重新登录！"}
	AuthExpireErr = ErrMsg{4002, "token 过期，请刷新token"}
	AuthActionErr = ErrMsg{4003, "权限错误"}
	ParamErr      = ErrMsg{4004, "参数解析失败，请联系管理员"}
	SystemErr     = ErrMsg{5000, "系统错误，请联系管理员"}
	DataEmptyErr  = ErrMsg{5001, "数据为空"}
	TokenCacheErr = ErrMsg{5002, "TOKEN CACHE 错误"}

	BizErrNameExist = ErrMsg{10100, "biz.err.name_exist"}
)

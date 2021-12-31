package httpUtils

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type Request struct {
	PageSize int `json:"pageSize"`
	PageNo   int `json:"pageNo"`
	Aa       int `json:"aa"`
}

type Response struct {
	Code    int         `json:"code"`
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
	Result  interface{} `json:"result"`
}
type PageResult struct {
	Items interface{} `json:"items"`
	Total int64       `json:"totalPage"`
}

func ApiRes(code consts.ResultCode, msg string, objects interface{}) (r *Response) {
	r = &Response{Code: code.Int(), Message: msg, Result: objects}
	return
}
func ApiResPage(code consts.ResultCode, msg string, objects interface{}, pageNo, pageSize int, total int64) (r *Response) {
	result := PageResult{Total: total, Items: objects}
	r = &Response{Code: code.Int(), Message: msg, Result: result}
	return
}

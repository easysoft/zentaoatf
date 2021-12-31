package serverDomain

import "github.com/aaronchen2k/deeptest/internal/pkg/consts"

type RpcResult struct {
	Code    consts.ResultCode `json:"code"`
	Msg     string            `json:"msg"`
	Payload interface{}       `json:"payload"`
}

func (result *RpcResult) Pass(msg string) {
	result.Code = consts.ResultCodeSuccess
	result.Msg = msg
}

func (result *RpcResult) Fail(msg string) {
	result.Code = consts.ResultCodeFail
	result.Msg = msg
}

func (result *RpcResult) IsSuccess() bool {
	return result.Code == consts.ResultCodeSuccess
}

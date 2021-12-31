package domain

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
)

type RpcReq struct {
	NodeIp   string
	NodePort int

	ApiPath   string
	ApiMethod string
	Data      interface{}
}

type RpcResp struct {
	Code    consts.ResultCode `json:"code"`
	Msg     string            `json:"msg"`
	Payload interface{}       `json:"payload"`
}

func (result *RpcResp) Pass(msg string) {
	result.Code = consts.ResultCodeSuccess
	result.Msg = msg
}

func (result *RpcResp) Passf(str string, args ...interface{}) {
	result.Code = consts.ResultCodeSuccess
	result.Msg = fmt.Sprintf(str, args...)
}

func (result *RpcResp) Fail(msg string) {
	result.Code = consts.ResultCodeFail
	result.Msg = msg
}

func (result *RpcResp) Failf(str string, args ...interface{}) {
	result.Code = consts.ResultCodeFail
	result.Msg = fmt.Sprintf(str, args...)
}

func (result *RpcResp) IsSuccess() bool {
	return result.Code == consts.ResultCodeSuccess
}

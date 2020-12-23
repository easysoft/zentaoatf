package domain

import serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"

type OptResult struct {
	Code    int         `json:"code,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func (result *OptResult) Success(msg string) {
	result.Code = serverConst.ResultSuccess.Int()
	result.Msg = msg
}

func (result *OptResult) Fail(msg string) {
	result.Code = serverConst.ResultFail.Int()
	result.Msg = msg
}

func (result *OptResult) IsSuccess() bool {
	return result.Code == serverConst.ResultSuccess.Int()
}

package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/pkg/domain"
)

type BaseCtrl struct {
}

func NewBaseCtrl() *BaseCtrl {
	return &BaseCtrl{}
}

func (c *BaseCtrl) SuccessResp(data interface{}) (ret domain.Response) {
	ret = domain.Response{Code: commConsts.Success.Code, Data: data}

	return
}

func (c *BaseCtrl) ErrResp(err commConsts.ResponseCode, msg string) (ret domain.Response) {
	ret = domain.Response{Code: err.Code, Msg: c.ErrMsg(err, msg)}

	return
}

func (c *BaseCtrl) BizErrResp(err *domain.BizError, msg string) (ret domain.Response) {
	ret = domain.Response{Code: err.Code, Msg: msg}

	return
}

func (c *BaseCtrl) ErrMsg(err commConsts.ResponseCode, msg string) (ret string) {
	//ret = i118Utils.Sprintf(err.Message)
	//
	//if ret != "" {
	//	ret += ": "
	//}

	ret += msg

	return
}

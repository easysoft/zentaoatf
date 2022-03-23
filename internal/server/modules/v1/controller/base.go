package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
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

func (c *BaseCtrl) ErrMsg(err commConsts.ResponseCode, msg string) (ret string) {
	//ret = i118Utils.Sprintf(err.Message)
	//
	//if ret != "" {
	//	ret += ": "
	//}

	ret += msg

	return
}

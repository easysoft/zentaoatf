package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
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
	ret = i118Utils.Sprintf(err.Key)

	if ret != "" {
		ret += ": "
	}

	ret += msg

	return
}

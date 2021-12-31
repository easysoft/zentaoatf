package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/str"
	"go.uber.org/zap"
)

type DataCtrl struct {
	DataService *service.DataService `inject:""`
}

func NewDataCtrl() *DataCtrl {
	return &DataCtrl{}
}

// InitDB 初始化项目接口
func (c *DataCtrl) Init(ctx iris.Context) {
	req := serverDomain.DataRequest{}
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	err := c.DataService.InitDB(req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: domain.SystemErr.Msg})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

// Check 检测是否需要初始化项目
func (c *DataCtrl) Check(ctx iris.Context) {
	if c.DataService.DataRepo.DB == nil {
		ctx.JSON(domain.Response{Code: domain.NeedInitErr.Code, Data: iris.Map{
			"needInit": true,
		}, Msg: str.Join(domain.NeedInitErr.Msg, ":数据库初始化失败")})
		return
	} else if serverConsts.CONFIG.System.CacheType == "redis" && serverConsts.CACHE == nil {
		ctx.JSON(domain.Response{Code: domain.NeedInitErr.Code, Data: iris.Map{
			"needInit": true,
		}, Msg: str.Join(domain.NeedInitErr.Msg, ":缓存驱动初始化失败")})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: iris.Map{
		"needInit": false,
	}, Msg: domain.NoErr.Msg})
}

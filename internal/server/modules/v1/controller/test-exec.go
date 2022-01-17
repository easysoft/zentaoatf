package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/exec"
	reportUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/report"
	"strings"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type TestExecCtrl struct {
	TestExecService *service.TestExecService `inject:""`
	BaseCtrl
}

func NewTestExecCtrl() *TestExecCtrl {
	return &TestExecCtrl{}
}

// Query 分页列表
func (c *TestExecCtrl) List(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	var req serverDomain.TestExecReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		logUtils.Errorf("参数验证失败", err.Error())
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	req.ConvertParams()

	data, err := c.TestExecService.Paginate(req, projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

// Get 详情
func (c *TestExecCtrl) Get(ctx iris.Context) {
	var reqId domain.ReqId
	if err := ctx.ReadParams(&reqId); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	execution, err := c.TestExecService.FindById(reqId.Id)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: domain.SystemErr.Msg})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: execution, Msg: domain.NoErr.Msg})
}

// ExecCase 添加
func (c *TestExecCtrl) ExecCase(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")
	req := serverDomain.TestExec{}

	if err := ctx.ReadJSON(&req); err != nil {
		logUtils.Errorf("参数验证失败", err.Error())
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	report, pathMaxWidth, _ := scriptUtils.ExecCase(req, projectPath)
	reportUtils.GenZTFTestReport(report, pathMaxWidth, projectPath)

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

// Update 更新
func (c *TestExecCtrl) Update(ctx iris.Context) {
	var reqId domain.ReqId
	if err := ctx.ReadParams(&reqId); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	var req model.TestExec
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err := c.TestExecService.Update(reqId.Id, req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

// Delete 删除
func (c *TestExecCtrl) Delete(ctx iris.Context) {
	var req domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}
	err := c.TestExecService.DeleteById(req.Id)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

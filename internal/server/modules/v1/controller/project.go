package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"strings"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type ProjectCtrl struct {
	ProjectService *service.ProjectService `inject:""`
	BaseCtrl
}

func NewProjectCtrl() *ProjectCtrl {
	return &ProjectCtrl{}
}

// Query 分页列表
func (c *ProjectCtrl) List(ctx iris.Context) {
	var req serverDomain.ProjectReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	req.ConvertParams()

	data, err := c.ProjectService.Paginate(req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

// Get 详情
func (c *ProjectCtrl) Get(ctx iris.Context) {
	var req domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}
	product, err := c.ProjectService.FindById(req.Id)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: domain.SystemErr.Msg})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: product, Msg: domain.NoErr.Msg})
}

// Create 添加
func (c *ProjectCtrl) Create(ctx iris.Context) {
	req := model.Project{}
	err := ctx.ReadJSON(&req)
	req.Path = strings.TrimSpace(req.Path)

	if err != nil || req.Path == "" {
		logUtils.Errorf("参数验证失败 %s", err.Error())
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: "参数验证失败"})
		return
	}

	id, err := c.ProjectService.Create(req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.CommonErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: iris.Map{"id": id}, Msg: domain.NoErr.Msg})
}

// Update 更新
func (c *ProjectCtrl) Update(ctx iris.Context) {
	var reqId domain.ReqId
	if err := ctx.ReadParams(&reqId); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	var req model.Project
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err := c.ProjectService.Update(reqId.Id, req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

// Delete 删除
func (c *ProjectCtrl) Delete(ctx iris.Context) {
	var req domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}
	err := c.ProjectService.DeleteById(req.Id)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

func (c *ProjectCtrl) GetByUser(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	if projectPath == "" {
		projectPath = commConsts.WorkDir
	}

	projects, currProject, currProjectConfig, scriptTree, err := c.ProjectService.GetByUser(projectPath)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ret := iris.Map{"projects": projects, "currProject": currProject,
		"currConfig": currProjectConfig, "scriptTree": scriptTree}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: ret, Msg: domain.NoErr.Msg})
}

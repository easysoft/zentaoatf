package controller

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
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
	err := ctx.ReadQuery(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}
	req.ConvertParams()

	data, err := c.ProjectService.Paginate(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

// Create 添加
func (c *ProjectCtrl) Create(ctx iris.Context) {
	req := model.Project{}
	err := ctx.ReadJSON(&req)
	if err != nil || req.Path == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
		return
	}

	_, err = c.ProjectService.Create(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

// Delete 删除
func (c *ProjectCtrl) Delete(ctx iris.Context) {
	projectPath := ctx.URLParam("path")

	if projectPath == "" {
		ctx.JSON(c.ErrResp(commConsts.ParamErr, "path"))
		return
	}
	err := c.ProjectService.DeleteByPath(projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *ProjectCtrl) GetByUser(ctx iris.Context) {
	projectPath := ctx.URLParam("currProject")

	if projectPath == "" {
		projectPath = commConsts.WorkDir
	}

	projects, currProject, currProjectConfig, scriptTree, err := c.ProjectService.GetByUser(projectPath)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.Failure, err.Error()))
		return
	}

	data := iris.Map{"projects": projects, "currProject": currProject,
		"currConfig": currProjectConfig, "scriptTree": scriptTree}

	ctx.JSON(c.SuccessResp(data))
}

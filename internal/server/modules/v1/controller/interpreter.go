package controller

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
)

type InterpreterCtrl struct {
	InterpreterService *service.InterpreterService `inject:""`
	WorkspaceService   *service.WorkspaceService   `inject:""`
	BaseCtrl
}

func NewInterpreterCtrl() *InterpreterCtrl {
	return &InterpreterCtrl{}
}

func (c *InterpreterCtrl) GetLangSettings(ctx iris.Context) {
	data, err := c.InterpreterService.GetLangSettings()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *InterpreterCtrl) GetLangInterpreter(ctx iris.Context) {
	language := ctx.URLParam("language")

	data, err := execHelper.GetLangInterpreter(language)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *InterpreterCtrl) List(ctx iris.Context) {
	data, err := c.InterpreterService.List()
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(data))
}

func (c *InterpreterCtrl) Get(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	po, err := c.InterpreterService.Get(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}
	ctx.JSON(c.SuccessResp(po))
}

func (c *InterpreterCtrl) Create(ctx iris.Context) {
	req := model.Interpreter{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	id, err := c.InterpreterService.Create(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrRecordWithSameInterpreterExist, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(iris.Map{"id": id}))
}

func (c *InterpreterCtrl) Update(ctx iris.Context) {
	req := model.Interpreter{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	}

	err := c.InterpreterService.Update(req)
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.ErrRecordWithSameInterpreterExist, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(iris.Map{"id": req.ID}))
}

func (c *InterpreterCtrl) Delete(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	err = c.InterpreterService.Delete(uint(id))
	if err != nil {
		ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
		return
	}

	c.WorkspaceService.UpdateAllConfig()

	ctx.JSON(c.SuccessResp(nil))
}

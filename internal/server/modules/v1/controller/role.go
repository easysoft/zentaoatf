package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/core/web/validate"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"strings"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type RoleCtrl struct {
	RoleService *service.RoleService `inject:""`
}

func NewRoleCtrl() *RoleCtrl {
	return &RoleCtrl{}
}

// GetAllRoles 分页列表
func (c *RoleCtrl) GetAllRoles(ctx iris.Context) {
	var req serverDomain.RoleReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	data, err := c.RoleService.Paginate(req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: data, Msg: domain.NoErr.Msg})
}

// GetRole 详情
func (c *RoleCtrl) GetRole(ctx iris.Context) {
	var req domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}
	role, err := c.RoleService.FindById(req.Id)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: domain.SystemErr.Msg})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: role, Msg: domain.NoErr.Msg})
}

// CreateRole 添加
func (c *RoleCtrl) CreateRole(ctx iris.Context) {
	req := serverDomain.RoleRequest{}
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	id, err := c.RoleService.Create(req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: iris.Map{"id": id}, Msg: domain.NoErr.Msg})
}

// UpdateRole 更新
func (c *RoleCtrl) UpdateRole(ctx iris.Context) {
	var reqId domain.ReqId
	if err := ctx.ReadParams(&reqId); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}

	var req serverDomain.RoleRequest
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			logUtils.Errorf("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err := c.RoleService.Update(reqId.Id, req)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}

// DeleteRole 删除
func (c *RoleCtrl) DeleteRole(ctx iris.Context) {
	var req domain.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		logUtils.Errorf("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(domain.Response{Code: domain.ParamErr.Code, Data: nil, Msg: domain.ParamErr.Msg})
		return
	}
	err := c.RoleService.DeleteById(req.Id)
	if err != nil {
		ctx.JSON(domain.Response{Code: domain.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(domain.Response{Code: domain.NoErr.Code, Data: nil, Msg: domain.NoErr.Msg})
}


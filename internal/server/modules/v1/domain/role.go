package serverDomain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
)

type RoleRequest struct {
	model.BaseRole
	Perms [][]string `json:"perms"`
}

type RoleReqPaginate struct {
	domain.PaginateReq
	Name string `json:"name"`
}

type RoleResponse struct {
	domain.Model
	model.BaseRole
}

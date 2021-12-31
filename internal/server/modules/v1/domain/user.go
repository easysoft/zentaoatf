package serverDomain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/snowlyg/helper/str"
	"regexp"
)

type UserRequest struct {
	model.BaseUser
	Password string `json:"password"`
	RoleIds  []uint `json:"role_ids"`
}

type UserReqPaginate struct {
	domain.PaginateReq
	Name string `json:"name"`
}

type UserResponse struct {
	domain.Model
	model.BaseUser
	Roles []string `gorm:"-" json:"roles"`
}

func (res *UserResponse) ToString() {
	if res.Avatar == "" {
		return
	}
	re := regexp.MustCompile("^http")
	if !re.MatchString(res.Avatar) {
		res.Avatar = str.Join("http://127.0.0.1:8085/upload/", res.Avatar)
	}
}

type LoginResponse struct {
	domain.ReqId
	Password string `json:"password"`
}

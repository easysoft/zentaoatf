package serverDomain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
)

type ProjectReqPaginate struct {
	domain.PaginateReq
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

type ProjectConfig struct {
	Version  string `json:"version"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

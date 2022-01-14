package serverDomain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
)

type TestExecReqPaginate struct {
	domain.PaginateReq
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

type TestExec struct {
	Name  string   `json:"name"`
	Cases []string `json:"cases"`
}

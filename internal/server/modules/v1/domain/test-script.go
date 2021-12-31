package serverDomain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
)

type TestScriptRequest struct {
	model.TestScript
}

type TestScriptReqPaginate struct {
	domain.PaginateReq
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

type TestScriptResponse struct {
	model.TestScript
}

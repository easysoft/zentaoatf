package serverDomain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
)

type TestCaseRequest struct {
	model.TestCase
}

type TestCaseReqPaginate struct {
	domain.PaginateReq
	Name     string `json:"name"`
	Category string `json:"name"`
	Status   string `json:"status"`
}

type TestCaseResponse struct {
	model.TestCase
}

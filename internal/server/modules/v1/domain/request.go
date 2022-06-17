package serverDomain

import (
	"github.com/easysoft/zentaoatf/pkg/domain"
)

type ReqPaginate struct {
	domain.PaginateReq
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

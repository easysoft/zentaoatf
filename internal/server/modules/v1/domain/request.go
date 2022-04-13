package serverDomain

import (
	"github.com/easysoft/zentaoatf/internal/pkg/domain"
)

type ReqPaginate struct {
	domain.PaginateReq
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

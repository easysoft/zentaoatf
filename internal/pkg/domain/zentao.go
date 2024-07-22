package commDomain

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/kataras/iris/v12"
)

type ZentaoUserProfile struct {
	Id       int    `json:"id"`
	Account  string `json:"account"`
	Realname string `json:"realname"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type ZentaoCaseStep struct {
	Type   commConsts.CaseStepType `json:"type"`
	Name   string                  `json:"name"`
	Desc   string                  `json:"desc"`
	Expect string                  `json:"expect"`
}

type BugOptionsWrapper struct {
	Options BugOptions `json:"options"`
}

type BugOptions struct {
	Type iris.Map      `json:"type"`
	Pri  []interface{} `json:"pri"`

	SeverityObj interface{} `json:"severity"`

	Modules iris.Map `json:"modules"`
	Build   iris.Map `json:"build"`
}

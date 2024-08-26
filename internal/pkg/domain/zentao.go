package commDomain

import (
	"github.com/kataras/iris/v12"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
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
	ID     string                  `json:"id"`
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

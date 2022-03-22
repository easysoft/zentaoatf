package commDomain

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type ZentaoUserProfile struct {
	Id       int    `json:"id"`
	Account  string `json:"account"`
	Realname string `json:"realname"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type ZentaoCaseStep struct {
	Type   commConsts.CaseStepType
	Desc   string
	Expect string
}

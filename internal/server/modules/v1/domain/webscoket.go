package serverDomain

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type WsReq struct {
	Act commConsts.ExecCmd `json:"act"`

	Seq   string `json:"seq"`
	Scope string `json:"scope"`

	Cases       []string `json:"cases"`
	ProductId   string   `json:"productId"`
	ModuleId    string   `json:"moduleId"`
	SuiteId     string   `json:"suiteId"`
	TaskId      string   `json:"taskId"`
	ProjectPath string   `json:"projectPath"`
}

type WsResp struct {
	Msg       string                   `json:"msg"`
	IsRunning string                   `json:"isRunning,omitempty"`
	Category  commConsts.WsMsgCategory `json:"category"`
}

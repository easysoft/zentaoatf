package serverDomain

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type WsReq struct {
	Act commConsts.ExecCmd `json:"act"`

	Seq   string                  `json:"seq"`
	Scope commConsts.ResultStatus `json:"scope"`

	Cases       []string `json:"cases"`
	ProductId   string   `json:"productId"`
	ModuleId    string   `json:"moduleId"`
	SuiteId     string   `json:"suiteId"`
	TaskId      string   `json:"taskId"`
	ProjectPath string   `json:"projectPath"`

	// for no-ztf testing like unittest etc.
	Framework    commConsts.UnitTestFramework `json:"framework"`
	Tool         commConsts.UnitTestFramework `json:"tool"`
	Cmd          string                       `json:"cmd"`
	SubmitResult bool                         `json:"submitResult"`
}

type WsResp struct {
	Msg       string                   `json:"msg"`
	IsRunning string                   `json:"isRunning,omitempty"`
	Category  commConsts.WsMsgCategory `json:"category"`
}

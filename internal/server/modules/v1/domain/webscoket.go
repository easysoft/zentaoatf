package serverDomain

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type WsReq struct {
	Act         commConsts.ExecCmd `json:"act"`
	Cases       []string           `json:"cases"`
	ModuleId    int                `json:"moduleId"`
	SuiteId     int                `json:"suiteId"`
	TaskId      int                `json:"taskId"`
	ProjectPath string             `json:"projectPath"`
}

type WsResp struct {
	Msg       string                   `json:"msg"`
	IsRunning string                   `json:"isRunning,omitempty"`
	Category  commConsts.WsMsgCategory `json:"category"`
}

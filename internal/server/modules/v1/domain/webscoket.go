package serverDomain

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type WsReq struct {
	Act commConsts.ExecCmd `json:"act"`

	Seq   string                  `json:"seq"`
	Scope commConsts.ResultStatus `json:"scope"`

	// for ztf testing
	Cases                     []string `json:"cases"`
	ProductId                 string   `json:"productId"`
	ModuleId                  string   `json:"moduleId"`
	SuiteId                   string   `json:"suiteId"`
	TaskId                    string   `json:"taskId"`
	WorkspaceId               int      `json:"workspaceId"`
	WorkspacePath             string   `json:"workspacePath"`
	ScriptDirParamFromCmdLine string   `json:"-"`

	// for unit, automation testing
	TestTool  commConsts.TestTool  `json:"testTool"`
	BuildTool commConsts.BuildTool `json:"buildTool"`
	Cmd       string               `json:"cmd"`

	SubmitResult bool `json:"submitResult"`
}

type WsResp struct {
	Msg       string                   `json:"msg"`
	IsRunning string                   `json:"isRunning,omitempty"`
	Category  commConsts.WsMsgCategory `json:"category"`
}

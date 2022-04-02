package serverDomain

import commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"

type WsReq struct {
	Act commConsts.ExecCmd `json:"act"`

	Seq   string                  `json:"seq"`
	Scope commConsts.ResultStatus `json:"scope"`

	// for ztf testing
	TestSets []TestSet `json:"testSets"`

	ProductId                 int    `json:"productId"`
	ModuleId                  int    `json:"moduleId"`
	SuiteId                   int    `json:"suiteId"`
	TaskId                    int    `json:"taskId"`
	ScriptDirParamFromCmdLine string `json:"-"`

	// for unit, automation testing
	TestTool  commConsts.TestTool  `json:"testTool"`
	BuildTool commConsts.BuildTool `json:"buildTool"`
	Cmd       string               `json:"cmd"`

	SubmitResult bool `json:"submitResult"`
}

type TestSet struct {
	WorkspaceId   int                 `json:"workspaceId"`
	WorkspaceType commConsts.TestTool `json:"workspaceType"`
	WorkspacePath string              `json:"workspacePath"`
	Cases         []string            `json:"cases"`

	Seq                       string                  `json:"seq"`
	Scope                     commConsts.ResultStatus `json:"scope"`
	ProductId                 int                     `json:"productId"`
	ModuleId                  int                     `json:"moduleId"`
	SuiteId                   int                     `json:"suiteId"`
	TaskId                    int                     `json:"taskId"`
	ScriptDirParamFromCmdLine string                  `json:"-"`

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

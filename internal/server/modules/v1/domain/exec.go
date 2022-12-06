package serverDomain

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
)

type ZentaoExecReq struct {
	Workspace string `json:"workspace"`
	Path      string `json:"path"`
	Ids       string `json:"ids"`
	Cmd       string `json:"cmd"`

	Task int `json:"task"`
}
type ZentaoCancelReq struct {
	Task int `json:"task"`
}

type ExecReq struct {
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
	Name string `json:"name"`

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
	ProxyId                   int                     `json:"proxyId"`

	// for unit, automation testing
	TestTool  commConsts.TestTool  `json:"testTool"`
	BuildTool commConsts.BuildTool `json:"buildTool"`
	Cmd       string               `json:"cmd"`

	ResultDir string `json:"resultDir"`
	ZipDir    string `json:"zipDir"`

	SubmitResult bool `json:"submitResult"`
}

type TestReportSummary struct {
	Name      string               `json:"name"`
	No        string               `json:"no"`
	Seq       string               `json:"seq"`
	TestEnv   commConsts.OsType    `json:"testEnv,omitempty"`
	TestType  commConsts.TestType  `json:"testType"`
	TestTool  commConsts.TestTool  `json:"testTool"`
	BuildTool commConsts.BuildTool `json:"buildTool"`

	ProductId     int               `json:"productId,omitempty"`
	WorkspaceId   int               `json:"workspaceId,omitempty"`
	WorkspaceName string            `json:"workspaceName,omitempty"`
	ExecBy        commConsts.ExecBy `json:"execBy,omitempty"`
	ExecById      int               `json:"execById,omitempty"`

	TestScriptName string `json:"testScriptName,omitempty"`

	Pass      int   `json:"pass"`
	Fail      int   `json:"fail"`
	Skip      int   `json:"skip"`
	Total     int   `json:"total"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	Duration  int64 `json:"duration"`
}

package serverDomain

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
)

type TestExecReq struct {
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

type TestReportSummary struct {
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

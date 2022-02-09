package serverDomain

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
)

type TestExecReq struct {
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

type TestReportSummary struct {
	Seq       string               `json:"seq"`
	TestEnv   commConsts.OsType    `json:"testEnv,omitempty"`
	TestType  commConsts.TestType  `json:"testType"`
	TestTool  commConsts.TestTool  `json:"testTool"`
	BuildTool commConsts.BuildTool `json:"buildTool"`

	ProductId int               `json:"productId,omitempty"`
	ExecBy    commConsts.ExecBy `json:"execBy,omitempty"`
	ExecById  int               `json:"execById,omitempty"`

	Pass      int   `json:"pass"`
	Fail      int   `json:"fail"`
	Skip      int   `json:"skip"`
	Total     int   `json:"total"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	Duration  int64 `json:"duration"`
}

package serverDomain

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
)

type TestExecReq struct {
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

type TestReportSummary struct {
	Seq      string `json:"seq"`
	Env      string `json:"env,omitempty"`
	TestType string `json:"testType"`

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

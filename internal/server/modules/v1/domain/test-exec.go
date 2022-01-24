package serverDomain

type TestExecReq struct {
	Keywords string `json:"keywords"`
	Enabled  string `json:"enabled"`
}

type TestReportSummary struct {
	Name     string `json:"name"`
	Env      string `json:"env,omitempty"`
	TestType string `json:"testType"`

	ProductId int `json:"productId,omitempty"`
	TaskId    int `json:"taskId,omitempty"`

	Pass      int   `json:"pass"`
	Fail      int   `json:"fail"`
	Skip      int   `json:"skip"`
	Total     int   `json:"total"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	Duration  int64 `json:"duration"`
}

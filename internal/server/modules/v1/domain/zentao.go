package serverDomain

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
)

type ZentaoResp struct {
	Status string
	Data   string
}
type ZentaoRespData struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type ZentaoJobSubmitReq struct {
	Task      int                  `json:"task"`
	Status    commConsts.JobStatus `json:"status"`
	StartTime string               `json:"startTime"`
	EndTime   string               `json:"endTime"`
	RetryTime int                  `json:"retryTime"`
	Error     string               `json:"error"`
	Data      interface{}          `json:"data"`
}

type ZentaoResultSubmitReq struct {
	Name        string `json:"name"`
	Seq         string `json:"seq"`
	WorkspaceId int    `json:"workspaceId"`
	ProductId   int    `json:"productId"`
	TaskId      int    `json:"taskId"`

	Task int `json:"task"`
}

type ZentaoLang struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ZentaoSite struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`

	Checked bool `json:"checked"`
}
type ZentaoProduct struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}

type ZentaoModule struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ZentaoSuite struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ZentaoTask struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

package serverDomain

type JobReq struct {
	JobId uint `json:"jobId"`

	Name      string `json:"name"`
	CaseIds   []uint `json:"caseIds"`
	ProductId int    `json:"productId"`
	ModuleId  int    `json:"moduleId"`
	SuiteId   int    `json:"suiteId"`
	TaskId    int    `json:"taskId"`
}

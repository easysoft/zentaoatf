package commDomain

type SyncSettings struct {
	WorkspaceId int `json:"workspaceId"`

	ProductId int    `json:"productId"`
	ModuleId  int    `json:"moduleId"`
	SuiteId   int    `json:"suiteId"`
	TaskId    int    `json:"taskId"`
	CaseId    int    `json:"caseId"`
	Lang      string `json:"lang"`

	SaveByModule    bool `json:"saveByModule"`
	IndependentFile bool `json:"independentFile"`
}

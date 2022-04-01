package serverDomain

type ZentaoResp struct {
	Status string
	Data   string
}
type ZentaoRespData struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type ZentaoResultSubmitReq struct {
	Seq         string `json:"seq"`
	WorkspaceId int    `json:"workspaceId"`
	ProductId   int    `json:"productId"`
	TaskId      int    `json:"taskId"`
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
}
type ZentaoProduct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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

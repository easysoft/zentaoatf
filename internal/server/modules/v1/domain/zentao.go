package serverDomain

type ZentaoResp struct {
	Status string
	Data   string
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

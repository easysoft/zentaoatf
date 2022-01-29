package commDomain

type ZentaoBugFields struct {
	Modules    []BugOption `json:"modules"`
	Categories []BugOption `json:"categories"`
	Versions   []BugOption `json:"versions"`
	Severities []BugOption `json:"severities"`
	Priorities []BugOption `json:"priorities"`
}

type BugOption struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

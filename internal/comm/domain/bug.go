package commDomain

type ZentaoBugFields struct {
	Types    []BugOption `json:"type"`
	Pri      []BugOption `json:"pri"`
	Severity []BugOption `json:"severity"`
	Modules  []BugOption `json:"modules"`
	Build    []BugOption `json:"build"`
}

type BugOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

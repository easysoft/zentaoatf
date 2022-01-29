package commDomain

type ZentaoBugFields struct {
	Modules    []BugOption
	Categories []BugOption
	Versions   []BugOption
	Severities []BugOption
	Priorities []BugOption
}

type BugOption struct {
	Id   string
	Code string
	Name string
}

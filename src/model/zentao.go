package model

type ZentaoResponse struct {
	Status string
	Data   string
}

type ZentaoBugFields struct {
	Modules    []Option
	Categories []Option
	Versions   []Option
	Severities []Option
	Priorities []Option
}

//type ZentaoCaseFileds struct {
//	Modules map[string]string
//}

type Option struct {
	Id   string
	Code string
	Name string
}

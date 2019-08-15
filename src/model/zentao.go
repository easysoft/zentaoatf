package model

type ZentaoResponse struct {
	Status string
	Data   string
}

type ZentaoSettings struct {
	Modules    []Option
	Categories []Option
	Versions   []Option
	Severities []Option
	Priorities []Option
}

type Option struct {
	Id        string
	Code      string
	Name      string
	IsDefault bool
}

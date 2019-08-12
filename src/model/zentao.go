package model

type ZentaoResponse struct {
	Status string
	Data   string
}

type ZentaoSettings struct {
	Modules    []Option
	Categories []Option
	Versions   []Option
	Priorities []Option
}

type Option struct {
	Id        int
	Code      string
	Name      string
	IsDefault bool
}

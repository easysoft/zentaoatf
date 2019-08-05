package model

type ZentaoSettings struct {
	Modules    []Option
	Categories []Option
	Versions   []Option
	Priorities []Option
}

type Option struct {
	Id        string
	Code      string
	Name      string
	isDefault bool
}

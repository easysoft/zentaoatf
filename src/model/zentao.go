package model

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
	isDefault bool
}

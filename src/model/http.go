package model

type Response struct {
	Code int
	Name string

	Modules    []Module
	Categories []Category
	Versions   []Version
	Priorities []Priority

	Cases []TestCase
}

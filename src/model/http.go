package model

type Response struct {
	Code int
	Name string

	ZentaoSettings ZentaoSettings

	Cases []TestCase
}

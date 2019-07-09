package model

type Response struct {
	Code  int
	Cases []TestCase
}

type TestCase struct {
	Id    int
	Title string
	Steps []TestStep
}

type TestStep struct {
	TestCase

	Expect       string
	IsGroup      bool
	IsCheckPoint bool
}

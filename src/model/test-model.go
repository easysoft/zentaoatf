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

type TestReport struct {
	Path string
	Env  string

	Pass      int
	Fail      int
	Total     int
	StartTime int64
	EndTime   int64
	Duration  int64

	Cases []CaseLog
}
type CaseLog struct {
	Numb   int
	Path   string
	Status bool

	CheckPoints []CheckPointLog
}
type CheckPointLog struct {
	Numb   int
	Expect string
	Actual string
	Status bool
}

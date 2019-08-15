package model

type Product struct {
	Id   int
	Code string
	Name string
}

type Module struct {
	Id   int
	Code string
	Name string
}

type TestTask struct {
	Id   int
	Code string
	Name string

	Runs map[int]TestCaseInTask
}

type TestCaseInTask struct {
	Id    string // runId in task
	Title string
	Case  string // real caseId
}

type TestCase struct {
	Id       string // runId in task
	IdInTask string

	Title   string
	Steps   map[int]TestStep
	StepArr []TestStep
}

type TestStep struct {
	Id   string
	Desc string

	Expect string
	Type   string
	Parent string
}

type TestReport struct {
	Path string
	Env  string

	Pass      int
	Fail      int
	Skip      int
	Total     int
	StartTime int64
	EndTime   int64
	Duration  int64

	Cases []CaseLog
}
type CaseLog struct {
	Id          int
	IdInTask    int
	ZentaoRunId int
	Path        string
	Status      string

	Steps []StepLog
}
type StepLog struct {
	Id     int
	Name   string
	Status bool

	CheckPoints []CheckPointLog
}
type CheckPointLog struct {
	Numb   int
	Expect string
	Actual string
	Status bool
}

type CaseResult struct {
	Case    string
	Version string

	Steps map[string]string
	Reals map[string]string
}

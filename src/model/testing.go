package model

type Product struct {
	Id   string
	Code string
	Name string
}

type Module struct {
	Id   int
	Code string
	Name string

	Cases map[int]TestCaseInModule
}
type TestCaseInModule struct {
	Id      string
	Title   string
	Product string
	Module  string
}

type TestSuite struct {
	Id      int
	Code    string
	Name    string
	Product string

	Cases map[int]TestCaseInSuite
}

type TestCaseInSuite struct {
	Id      string
	Title   string
	Product string
	Module  string
}

type TestTask struct {
	Id      int
	Code    string
	Name    string
	Product string
	Project string

	Runs map[int]TestCaseInTask
}

type TestCaseInTask struct {
	Id      string // runId in task
	Title   string
	Case    string // real caseId
	Product string
	Module  string
}

type TestCaseNoStepArr struct {
	Id      string
	Product string
	Module  string

	Title string
	Steps map[int]TestStep
}
type TestCase struct {
	Id      string
	Product string
	Module  string

	Title   string
	Steps   map[int]TestStep
	StepArr []TestStep `json tag -`
}

type TestStep struct {
	Id   string
	Desc string

	Expect string
	Type   string
	Parent string

	Children []TestStep
	Numb     string
}

type Bug struct {
	Title string

	Module      string            // id
	Type        string            // install
	OpenedBuild map[string]string // {"0": "trunk"}
	Severity    string            // id
	Pri         string            // id

	Product string
	Case    string
	Steps   string

	Uid         string // uuid.NewV4().String()
	CaseVersion string // 0
	OldTaskID   string // 0
}

type TestReport struct {
	Env string

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
	Id        int
	ProductId int
	Path      string
	Status    string
	Title     string

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

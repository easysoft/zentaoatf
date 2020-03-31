package model

import "encoding/xml"

type Product struct {
	Id   string
	Code string
	Name string

	Cases map[int]TestCaseInModule
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
type TestCaseWrapper struct {
	From string
	Case TestCase
}

type TestStep struct {
	Id   string
	Desc string

	Expect string
	Type   string
	Parent string

	Children []TestStep
	Numb     string

	MutiLine bool
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
	Env       string
	TestType  string
	TestFrame string

	Pass      int
	Fail      int
	Skip      int
	Total     int
	StartTime int64
	EndTime   int64
	Duration  int64

	ZtfCaseResults  []ZtfCaseResult
	UnitCaseResults []UnitCaseResult
}

type ZtfCaseResult struct {
	Id        int
	ProductId int
	Path      string
	Status    string
	Title     string

	Steps []StepLog
}
type StepLog struct {
	Id     string
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

// 单元测试
type UnitTestSuite struct {
	XMLName xml.Name `xml:"testsuite"`

	Name     string
	Duration int

	Properties Properties       `xml:"properties"`
	TestCases  []UnitCaseResult `xml:"testcase"`
}
type UnitCaseResult struct {
	Title     string   `xml:"name,attr"`
	TestSuite string   `xml:"classname,attr"`
	Duration  float32  `xml:"time,attr"`
	Failure   *Failure `xml:"failure,omitempty"`

	Id     int
	Status string
}

type Failure struct {
	Type string `xml:"type,attr"`
	Desc string `xml:",innerxml"`
}

type Properties struct {
	Property []Property `xml:"property"`
}
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// jtest xml
type JTestSuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	Title      string   `xml:"name,attr"`
	TestSuites []struct {
		Title     string           `xml:"name,attr"`
		TestCases []UnitCaseResult `xml:"testcase"`

		Duration int
	} `xml:"testsuite"`

	Duration int
}

// phpunit xml
type PhpUnitSuites struct {
	XMLName   xml.Name `xml:"tests"`
	TestCases []struct {
		Title     string `xml:"prettifiedMethodName,attr"`
		TestSuite string `xml:"prettifiedClassName,attr"`
		Fail      string `xml:"exceptionMessage,attr"`

		Groups   string  `xml:"groups,attr"`
		Status   int     `xml:"status,attr"`
		Duration float32 `xml:"time,attr"`
	} `xml:"test"`

	Duration int
}

// pytest xml
type PyTestSuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	TestSuites []struct {
		Title     string `xml:"name,attr"`
		TestCases []struct {
			Title     string  `xml:"name,attr"`
			TestSuite string  `xml:"classname,attr"`
			Duration  float32 `xml:"time,attr"`
			Failure   *struct {
				Type string `xml:"message,attr"`
				Desc string `xml:",innerxml"`
			} `xml:"failure,omitempty"`

			Status string
		} `xml:"testcase"`

		Duration int
	} `xml:"testsuite"`

	Duration int
}

// gtest xml
type GTestSuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	TestSuites []struct {
		Title     string `xml:"name,attr"`
		TestCases []struct {
			Title     string  `xml:"name,attr"`
			TestSuite string  `xml:"classname,attr"`
			Duration  float32 `xml:"time,attr"`
			Failure   *struct {
				Type string `xml:"message,attr"`
				Desc string `xml:",innerxml"`
			} `xml:"failure,omitempty"`

			Status string `xml:"status,attr"`
		} `xml:"testcase"`

		Duration int
	} `xml:"testsuite"`

	Duration int
}

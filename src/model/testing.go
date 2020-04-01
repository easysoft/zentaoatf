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
	Env       string `json:"env"`
	TestType  string `json:"testType"`
	TestFrame string `json:"TestFrame"`

	ProductId  int    `json:"ProductId"`
	TaskId     int    `json:"TaskId"`
	ZentaoData string `json:"ZentaoData"`
	BuildUrl   string `json:"BuildUrl"`

	Pass      int   `json:"Pass"`
	Fail      int   `json:"Fail"`
	Skip      int   `json:"Skip"`
	Total     int   `json:"Total"`
	StartTime int64 `json:"StartTime"`
	EndTime   int64 `json:"EndTime"`
	Duration  int64 `json:"Duration"`

	ZTFCaseResults  []ZTFCaseResult  `json:"ZTFCaseResults"`
	UnitCaseResults []UnitCaseResult `json:"UnitCaseResults"`
}

type ZTFCaseResult struct {
	Id        int    `json:"Id"`
	ProductId int    `json:"ProductId"`
	Path      string `json:"Path"`
	Status    string `json:"Status"`
	Title     string `json:"Title"`

	Steps []StepLog `json:"Steps"`
}
type StepLog struct {
	Id     string `json:"Id"`
	Name   string `json:"Name"`
	Status bool   `json:"Status"`

	CheckPoints []CheckPointLog `json:"CheckPoints"`
}
type CheckPointLog struct {
	Numb   int    `json:"Numb"`
	Expect string `json:"Expect"`
	Actual string `json:"Actual"`
	Status bool   `json:"Status"`
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
	Title     string   `json:"title" xml:"name,attr"`
	TestSuite string   `json:"testSuite" xml:"classname,attr"`
	Duration  float32  `json:"duration" xml:"time,attr"`
	Failure   *Failure `json:"failure" xml:"failure,omitempty"`

	Id     int    `json:"id"`
	Status string `json:"status"`
}

type Failure struct {
	Type string `json:"type" xml:"type,attr"`
	Desc string `json:"desc" xml:",innerxml"`
}

type Properties struct {
	Property []Property `json:"property" xml:"property"`
}
type Property struct {
	Name  string `json:"name" xml:"name,attr"`
	Value string `json:"value" xml:"value,attr"`
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

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
	TestFrame string `json:"testFrame"`

	ProductId  int    `json:"productId"`
	TaskId     int    `json:"taskId"`
	ZentaoData string `json:"zentaoData"`
	BuildUrl   string `json:"buildUrl"`

	Pass      int   `json:"pass"`
	Fail      int   `json:"fail"`
	Skip      int   `json:"skip"`
	Total     int   `json:"total"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	Duration  int64 `json:"duration"`

	ZTFResults  []ZTFResult  `json:"ztfResults"`
	UnitResults []UnitResult `json:"unitResults"`
}

type ZTFResult struct {
	Id        int    `json:"id"`
	ProductId int    `json:"productId"`
	Path      string `json:"path"`
	Status    string `json:"status"`
	Title     string `json:"title"`

	Steps []StepLog `json:"steps"`
}
type StepLog struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`

	CheckPoints []CheckPointLog `json:"checkPoints"`
}
type CheckPointLog struct {
	Numb   int    `json:"numb"`
	Expect string `json:"expect"`
	Actual string `json:"actual"`
	Status bool   `json:"status"`
}

// 单元测试
type UnitTestSuite struct {
	XMLName xml.Name `xml:"testsuite"`

	Name     string
	Duration int

	Properties Properties   `xml:"properties"`
	TestCases  []UnitResult `xml:"testcase"`
}
type UnitResult struct {
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
		Title     string       `xml:"name,attr"`
		TestCases []UnitResult `xml:"testcase"`

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

// cppunit xml
type CppUnitSuites struct {
	XMLName xml.Name `xml:"TestRun"`

	SuccessfulTests struct {
		TestCases []CppUnitTest `json:"test" xml:"Test"`
	} `json:"successfulTests" xml:"SuccessfulTests"`

	FailedTests struct {
		TestCases []CppUnitTest `json:"test" xml:"FailedTest"`
	} `json:"failedTests" xml:"FailedTests"`

	Duration int
}
type CppUnitTest struct {
	Id          int    `json:"id" xml:"Id,attr"`
	Title       string `json:"name" xml:"Name,innerxml"`
	FailureType string `json:"failureType" xml:"FailureType,innerxml"`
	Message     string `json:"message" xml:"Message,innerxml"`
	Location    []struct {
		File string `json:"file" xml:"File,innerxml"`
		Line string `json:"line" xml:"Line,innerxml"`
	} `json:"location" xml:"Location"`

	Duration int
}

// qtest xml
type QTestSuites struct {
	XMLName xml.Name `xml:"testsuite"`
	Name    string   `json:"name" xml:"name,attr"`

	TestCases []struct {
		Title  string `json:"name" xml:"name,attr"`
		Result string `json:"result" xml:"result,attr"`

		Failure *struct {
			Type string `json:"type" xml:"tag,attr"`
			Desc string `json:"desc" xml:"message,attr"`
		} `json:"failure" xml:"failure"`
	} `json:"testCases" xml:"testcase"`

	Properties Properties `json:"properties" xml:"properties"`
	Duration   int
}

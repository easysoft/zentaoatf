package commDomain

import (
	"encoding/xml"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"time"
)

type ZtfProduct struct {
	Id   string
	Code string
	Name string

	Cases map[int]ZtfCaseInModule
}

type ZtfModule struct {
	Id   string
	Code string
	Name string

	Cases map[int]ZtfCaseInModule
}
type ZtfCaseInModule struct {
	Id      string
	Title   string
	Product string
	Module  string
}

type ZtfSuite struct {
	Id      string
	Code    string
	Name    string
	Product string

	Cases map[int]ZtfCaseInSuite
}

type ZtfCaseInSuite struct {
	Id      string
	Title   string
	Product string
	Suite   string
	Module  string
}

type ZtfTask struct {
	Id      string
	Code    string
	Name    string
	Product string
	Project string

	Runs map[int]ZtfCaseInTask
}

type ZtfCaseInTask struct {
	Id      string // runId in task
	Title   string
	Case    string // real caseId
	Product string
	Module  string
}

type ZtfCaseNoStepArr struct {
	Id      string
	Product string
	Module  string

	Title string
	Steps map[int]ZtfStep
}
type ZtfCase struct {
	Id      string
	Product string
	Module  string

	Title   string
	Steps   map[int]ZtfStep
	StepArr []ZtfStep `json tag -`
}
type ZtfCaseWrapper struct {
	From string
	Case ZtfCase
}

type ZtfStep struct {
	Id   string
	Desc string

	Expect string
	Type   string
	Parent string

	Children []ZtfStep
	Numb     string

	MultiLine bool
}

type ZtfBug struct {
	Title string `json:"title"`

	Module      int               `json:"module"`
	Type        string            `json:"type"`
	OpenedBuild map[string]string `json:"openedBuild"`
	Severity    string            `json:"severity"`
	Pri         string            `json:"pri"`

	Product int    `json:"product"`
	Case    string `json:"case"`
	Steps   string `json:"steps"`

	Uid         string `json:"uid"`
	CaseVersion string `json:"caseVersion"`
	OldTaskID   string `json:"oldTaskID"`
}

type ZtfReport struct {
	TestEnv   consts.OsType   `json:"testEnv,omitempty"`
	TestType  consts.TestType `json:"testType"`
	TestFrame string          `json:"testFrame"`

	ProductId int           `json:"productId,omitempty"`
	ExecBy    consts.ExecBy `json:"execBy,omitempty"`
	ExecById  int           `json:"execById,omitempty"`

	Pass      int   `json:"pass"`
	Fail      int   `json:"fail"`
	Skip      int   `json:"skip"`
	Total     int   `json:"total"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	Duration  int64 `json:"duration"`

	FuncResult []FuncResult `json:"funcResult,omitempty"`
	UnitResult []UnitResult `json:"unitResult,omitempty"`
}

type FuncResult struct {
	Id        int                     `json:"id"`
	ProductId int                     `json:"productId"`
	Path      string                  `json:"path"`
	Status    commConsts.ResultStatus `json:"status"`
	Title     string                  `json:"title"`

	Steps []StepLog `json:"steps"`
}
type StepLog struct {
	Id     string                  `json:"id"`
	Name   string                  `json:"name"`
	Status commConsts.ResultStatus `json:"status"`

	CheckPoints []CheckPointLog `json:"checkPoints"`
}
type CheckPointLog struct {
	Numb   int                     `json:"numb"`
	Expect string                  `json:"expect"`
	Actual string                  `json:"actual"`
	Status commConsts.ResultStatus `json:"status"`
}

// 单元测试
type UnitTestSuite struct {
	XMLName xml.Name `xml:"testsuite"`

	Name     string
	Duration int64   `xml:"-"`
	Time     float32 `xml:"time,attr"`

	Properties Properties   `xml:"properties"`
	ZtfScripts []UnitResult `xml:"testcase"`
}
type UnitResult struct {
	Title     string `json:"title" xml:"name,attr"`
	TestSuite string `json:"testSuite" xml:"classname,attr"`

	StartTime int64 `json:"startTime" xml:"startTime"`
	EndTime   int64 `json:"endTime" xml:"endTime"`

	Duration float32  `json:"duration" xml:"time,attr"`
	Failure  *Failure `json:"failure" xml:"failure,omitempty"`

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

// phpunit xml
type PhpUnitSuites struct {
	XMLName    xml.Name `xml:"tests"`
	ZtfScripts []struct {
		Title     string `xml:"prettifiedMethodName,attr"`
		TestSuite string `xml:"prettifiedClassName,attr"`
		Fail      string `xml:"exceptionMessage,attr"`

		Groups string  `xml:"groups,attr"`
		Status int     `xml:"status,attr"`
		Time   float32 `xml:"time,attr"`
	} `xml:"test"`

	Duration int
}

// pytest xml
type PyTestSuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	TestSuites []struct {
		Title      string `xml:"name,attr"`
		ZtfScripts []struct {
			Title     string  `xml:"name,attr"`
			TestSuite string  `xml:"classname,attr"`
			Duration  float32 `xml:"time,attr"`
			Failure   *struct {
				Type string `xml:"message,attr"`
				Desc string `xml:",innerxml"`
			} `xml:"failure,omitempty"`
			Error *struct {
				Text    string `xml:",chardata"`
				Message string `xml:"message,attr"`
			} `xml:"error"`

			Status string
		} `xml:"testcase"`

		Duration int
		Time     float32 `xml:"time,attr"`
	} `xml:"testsuite"`

	Duration int
}

// jest xml
type JestSuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	Title      string   `xml:"name,attr"`
	TestSuites []struct {
		Title      string       `xml:"name,attr"`
		ZtfScripts []UnitResult `xml:"testcase"`

		Duration int
	} `xml:"testsuite"`

	Duration int
	Time     float32 `xml:"time,attr"`
}

// gtest xml
type GTestSuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	TestSuites []struct {
		Title      string `xml:"name,attr"`
		ZtfScripts []struct {
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
	Time     float32 `xml:"time,attr"`
}

// qtest xml
type QTestSuites struct {
	XMLName xml.Name `xml:"testsuite"`
	Name    string   `json:"name" xml:"name,attr"`

	ZtfScripts []struct {
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

// cppunit xml
type CppUnitSuites struct {
	XMLName xml.Name `xml:"TestRun"`

	SuccessfulTests struct {
		ZtfScripts []CppUnitTest `json:"test" xml:"Tests"`
	} `json:"successfulTests" xml:"SuccessfulTests"`

	FailedTests struct {
		ZtfScripts []CppUnitTest `json:"test" xml:"FailedTest"`
	} `json:"failedTests" xml:"FailedTests"`

	Duration int
}
type CppUnitTest struct {
	Id          int    `json:"id" xml:"Id,attr"`
	Title       string `json:"name" xml:"Name"`
	FailureType string `json:"failureType" xml:"FailureType"`
	Message     string `json:"message" xml:"Message"`
	Location    []struct {
		File string `json:"file" xml:"File"`
		Line string `json:"line" xml:"Line"`
	} `json:"location" xml:"Location"`

	Duration int
}

// RobotFramework xml
var RFResults = "results"

type RobotResult struct {
	XMLName    xml.Name        `xml:"robot"`
	Text       string          `xml:",chardata"`
	Generator  string          `xml:"generator,attr"`
	Generated  string          `xml:"generated,attr"`
	Rpa        string          `xml:"rpa,attr"`
	Suites     []RobotSuite    `xml:"suite"`
	Statistics RobotStatistics `xml:"statistics"`
	Errors     string          `xml:"errors"`
}

type RobotStatistics struct {
	Text  string `xml:",chardata"`
	Total struct {
		Text string `xml:",chardata"`
		Stat []struct {
			Text string `xml:",chardata"`
			Pass string `xml:"pass,attr"`
			Fail string `xml:"fail,attr"`
		} `xml:"stat"`
	} `xml:"total"`
	Tag   string `xml:"tag"`
	Suite struct {
		Text   string `xml:",chardata"`
		States []struct {
			Text string `xml:",chardata"`
			Pass string `xml:"pass,attr"`
			Fail string `xml:"fail,attr"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"stat"`
	} `xml:"suite"`
}

type RobotSuite struct {
	Text   string       `xml:",chardata"`
	ID     string       `xml:"id,attr"`
	Name   string       `xml:"name,attr"`
	Source string       `xml:"source,attr"`
	Suites []RobotSuite `xml:"suite"`
	Tests  []RobotTest  `xml:"test"`
	Status RobotStatus  `xml:"status"`
}

type RobotTest struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Kw   []struct {
		Text      string `xml:",chardata"`
		Name      string `xml:"name,attr"`
		Library   string `xml:"library,attr"`
		Doc       string `xml:"doc"`
		Arguments struct {
			Text string   `xml:",chardata"`
			Arg  []string `xml:"arg"`
		} `xml:"arguments"`
		Msg struct {
			Text      string `xml:",chardata"`
			Timestamp string `xml:"timestamp,attr"`
			Level     string `xml:"level,attr"`
		} `xml:"msg"`
		Status RobotStatus `xml:"status"`
		Assign struct {
			Text string `xml:",chardata"`
			Var  string `xml:"var"`
		} `xml:"assign"`
	} `xml:"kw"`
	Doc    string      `xml:"doc"`
	Status RobotStatus `xml:"status"`
}

type RobotStatus struct {
	Text      string `xml:",chardata"`
	Status    string `xml:"status,attr"`
	StartTime string `xml:"starttime,attr"`
	EndTime   string `xml:"endtime,attr"`
}

// cypress
var CypressResults = "results"

type CypressTestsuites struct {
	XMLName    xml.Name           `xml:"testsuites"`
	Text       string             `xml:",chardata"`
	Name       string             `xml:"name,attr"`
	Time       string             `xml:"time,attr"`
	Tests      string             `xml:"tests,attr"`
	Failures   string             `xml:"failures,attr"`
	Testsuites []CypressTestsuite `xml:"testsuite"`
}

type CypressTestsuite struct {
	Text      string            `xml:",chardata"`
	Name      string            `xml:"name,attr"`
	Timestamp string            `xml:"timestamp,attr"`
	Tests     string            `xml:"tests,attr"`
	File      string            `xml:"file,attr"`
	Time      float64           `xml:"time,attr"`
	Failures  string            `xml:"failures,attr"`
	Testcases []CypressTestcase `xml:"testcase"`
}

type CypressTestcase struct {
	Text      string           `xml:",chardata"`
	Name      string           `xml:"name,attr"`
	Time      float64          `xml:"time,attr"`
	Classname string           `xml:"classname,attr"`
	Failures  []CypressFailure `xml:"failure"`
}

type CypressFailure struct {
	Text    string `xml:",chardata"`
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

type TestResult struct {
	TestSetId uint `json:"testSetId" yaml:"testSetId"`

	Version float64 `json:"version" yaml:"version"`
	Name    string  `json:"name" yaml:"name"`
	Code    int     `json:"code"`
	Msg     string  `json:"msg"`

	StartTime time.Time `json:"startTime" yaml:"startTime"`
	EndTime   time.Time `json:"endTime" yaml:"endTime"`
	Duration  int       `json:"duration" yaml:"duration"` // sec

	TotalNum  int `json:"totalNum" yaml:"totalNum"`
	PassNum   int `json:"passNum" yaml:"passNum"`
	FailNum   int `json:"failNum" yaml:"failNum"`
	MissedNum int `json:"missedNum" yaml:"missedNum"`

	Payload interface{} `json:"payload"`
}

func (result *TestResult) Pass(msg string) {
	result.Code = consts.ResultCodeSuccess.Int()
	result.Msg = msg
}
func (result *TestResult) Passf(str string, args ...interface{}) {
	result.Code = consts.ResultCodeSuccess.Int()
	result.Msg = fmt.Sprintf(str+"\n", args...)
}

func (result *TestResult) Fail(msg string) {
	result.Code = consts.ResultCodeFail.Int()
	result.Msg = msg
}

func (result *TestResult) Failf(str string, args ...interface{}) {
	result.Code = consts.ResultCodeFail.Int()
	result.Msg = fmt.Sprintf(str+"\n", args...)
}

func (result *TestResult) IsSuccess() bool {
	return result.Code == consts.ResultCodeSuccess.Int()
}

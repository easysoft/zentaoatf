package commDomain

import (
	"encoding/xml"
	"fmt"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/pkg/consts"
)

type ZtfRespTestCases struct {
	Cases []ZtfCaseInModule `json:"testcases"`
}

type ZtfProduct struct {
	Id int
	//Code string
	//Name string

	Cases map[int]ZtfCaseInModule `json:"testcases"`
}

type ZtfModule struct {
	Id    int
	Code  string
	Title string

	Cases map[int]ZtfCaseInModule
}
type ZtfCaseInModule struct {
	Id            int
	Title         string
	Product       int
	Module        int
	Case          int // case id in task
	Type          string
	StatusName    string
	LastRunResult string
}

type ZtfSuite struct {
	Id      int
	Code    string
	Name    string
	Product int

	Cases map[int]ZtfCaseInSuite
}

type ZtfCaseInSuite struct {
	Id      int
	Title   string
	Product int
	Suite   int
	Module  int
}

type ZtfTask struct {
	Id        int
	Code      string
	Name      string
	Product   int
	Workspace int

	Runs map[int]ZtfCaseInTask
}

type ZtfCaseInTask struct {
	Id      int // runId in task
	Title   string
	Case    int // real caseId
	Product int
	Module  int
}

type ZtfCaseNoStepArr struct {
	Id      int
	Product int
	Module  int

	Title string
	Steps map[int]ZtfStep
}
type ZtfCase struct {
	Id      int
	Product int
	Module  int

	Title string
	Steps []ZtfStep `json:"steps"`

	ScriptPath string `json:"scriptPath"` // used for update exist script
}
type ZtfCaseWrapper struct {
	From string
	Case ZtfCase
}

type ZtfStep struct {
	Id   int
	Desc string

	Expect string
	Type   string
	Parent int

	Children []ZtfStep
	Numb     string

	MultiLine bool
}

type ZtfBug struct {
	Title    string              `json:"title"`
	Type     string              `json:"type"`
	TestType commConsts.TestType `json:"testType"`

	StepIds string `json:"ids"` // for to

	Product  int `json:"product"`
	Module   int `json:"module"`
	Severity int `json:"severity"`
	Pri      int `json:"pri"`

	Case  int    `json:"case"`
	Steps string `json:"steps"`

	Uid         string   `json:"uid"`
	OpenedBuild []string `json:"openedBuild"`
	CaseVersion string   `json:"caseVersion"`
	OldTaskID   string   `json:"oldTaskID"`
}

type ZentaoBug struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`

	Product  int `json:"product"`
	Module   int `json:"module"`
	Severity int `json:"severity"`
	Pri      int `json:"pri"`

	Case  int    `json:"case"`
	Steps string `json:"steps"`

	Uid         string   `json:"uid"`
	OpenedBuild []string `json:"openedBuild"`
	CaseVersion string   `json:"caseVersion"`
	OldTaskID   string   `json:"oldTaskID"`
	OpenedDate  string   `json:"openedDate"`
	OpenedBy    string   `json:"openedBy"`
	StatusName  string   `json:"statusName"`
}

type ZtfReport struct {
	Name string `json:"name"`

	Platform    commConsts.OsType    `json:"platform,omitempty"`
	TestType    commConsts.TestType  `json:"testType"`
	TestTool    commConsts.TestTool  `json:"testTool"`
	BuildTool   commConsts.BuildTool `json:"buildTool"`
	TestCommand string               `json:"testCommand"`

	WorkspaceId   int                 `json:"workspaceId,omitempty"`
	WorkspaceType commConsts.TestTool `json:"workspaceType"`
	WorkspacePath string              `json:"workspacePath"`
	Seq           string              `json:"seq,omitempty"`
	SubmitResult  bool                `json:"submitResult"`

	ProductId int               `json:"productId,omitempty"`
	TaskId    int               `json:"taskId,omitempty"`
	Task      int               `json:"task,omitempty"`
	ExecBy    commConsts.ExecBy `json:"execBy,omitempty"`
	ExecById  int               `json:"execById,omitempty"`

	// run with ci tool
	ZentaoData string `json:"zentaoData"`
	BuildUrl   string `json:"buildUrl"`
	Log        string `json:"log"`

	Pass      int   `json:"pass"`
	Fail      int   `json:"fail"`
	Skip      int   `json:"skip"`
	Total     int   `json:"total"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	Duration  int64 `json:"duration"`

	FuncResult []FuncResult `json:"funcResult,omitempty"`
	UnitResult []UnitResult `json:"unitResult,omitempty"`

	JacocoResult *JacocoResult `json:"jacocoResult,omitempty"`
	ZapResult    *ZapResult    `json:"zapResult,omitempty"`
}

type FuncResult struct {
	Id          int    `json:"id"`
	WorkspaceId int    `json:"workspaceId"`
	Seq         string `json:"seq"`

	Key          string                  `json:"key"`
	ProductId    int                     `json:"productId"`
	Path         string                  `json:"path"`
	RelativePath string                  `json:"relativePath"`
	Status       commConsts.ResultStatus `json:"status"`
	Title        string                  `json:"title"`

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
	Cases      []UnitResult `xml:"testcase"`
}
type UnitResult struct {
	Title     string `json:"title" xml:"name,attr"`
	TestSuite string `json:"testSuite" xml:"classname,attr"`

	StartTime int64 `json:"startTime" xml:"startTime"`
	EndTime   int64 `json:"endTime" xml:"endTime"`

	Duration float32  `json:"duration" xml:"time,attr"`
	Failure  *Failure `json:"failure" xml:"failure,omitempty"`

	ErrorType    string `json:"errorType" xml:"type,attr,omitempty"`
	ErrorContent string `json:"errorContent" xml:"error,omitempty"`

	Id     int                     `json:"id"`
	Cid    int                     `json:"cid"`
	Status commConsts.ResultStatus `json:"status"`
}
type CaseSlice []UnitResult

func (a CaseSlice) Len() int {
	return len(a)
}
func (a CaseSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a CaseSlice) Less(i, j int) bool {
	return a[j].Id > a[i].Id
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

// allure json
type AllureCase struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	StatusDetails struct {
		Known   bool   `json:"known"`
		Muted   bool   `json:"muted"`
		Flaky   bool   `json:"flaky"`
		Message string `json:"message"`
		Trace   string `json:"trace"`
	} `json:"statusDetails"`
	Stage       string        `json:"stage"`
	Description string        `json:"description"`
	Steps       []interface{} `json:"steps"`
	Attachments []interface{} `json:"attachments"`
	Parameters  []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"parameters"`
	Start      int64  `json:"start"`
	Stop       int64  `json:"stop"`
	Uuid       string `json:"uuid"`
	HistoryId  string `json:"historyId"`
	TestCaseId string `json:"testCaseId"`
	FullName   string `json:"fullName"`
	Labels     []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"labels"`
	Links []interface{} `json:"links"`

	Duration int
}

// phpunit xml
type PhpUnitSuites struct {
	XMLName xml.Name `xml:"tests"`
	Cases   []struct {
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
		Title string `xml:"name,attr"`
		Cases []struct {
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
		Title string       `xml:"name,attr"`
		Cases []UnitResult `xml:"testcase"`

		Duration int
	} `xml:"testsuite"`

	Duration int
	Time     float32 `xml:"time,attr"`
}

// gtest xml
type GTestSuites struct {
	XMLName    xml.Name `xml:"testsuites"`
	TestSuites []struct {
		Title string `xml:"name,attr"`
		Cases []struct {
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

	Cases []struct {
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
		Cases []CppUnitTest `json:"test" xml:"Tests"`
	} `json:"successfulTests" xml:"SuccessfulTests"`

	FailedTests struct {
		Cases []CppUnitTest `json:"test" xml:"FailedTest"`
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

	OutPut string `xml:"system-out"`
}

type CypressFailure struct {
	Text    string `xml:",chardata"`
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

// k6
type K6Metric struct {
	Type   commConsts.K6ResultType `json:"type"`
	Metric string                  `json:"metric"`
	Data   struct {
		Name       string        `json:"name"`
		Type       string        `json:"type"`
		Contains   string        `json:"contains"`
		Thresholds []interface{} `json:"thresholds"`
		Submetrics interface{}   `json:"submetrics"`
	} `json:"data"`
}
type K6Point struct {
	Type   commConsts.K6ResultType `json:"type"`
	Metric string                  `json:"metric"`
	Data   struct {
		Time  time.Time `json:"time"`
		Value int       `json:"value"`
		Tags  struct {
			ExpectedResponse string `json:"expected_response"`
			Group            string `json:"group"`
			Method           string `json:"method"`
			Proto            string `json:"proto"`
			Scenario         string `json:"scenario"`
			Status           string `json:"status"`
			TlsVersion       string `json:"tls_version"`
			Url              string `json:"url"`

			Id         string `json:"id"`
			Name       string `json:"name"`
			Checkpoint string `json:"checkpoint"`
		} `json:"tags"`
	} `json:"data"`
}
type K6Summary struct {
	RootGroup struct {
		Name   string `json:"name"`
		Path   string `json:"path"`
		Id     string `json:"id"`
		Groups []struct {
			Name   string        `json:"name"`
			Path   string        `json:"path"`
			Id     string        `json:"id"`
			Groups []interface{} `json:"groups"`
			Checks []struct {
				Path   string `json:"path"`
				Id     string `json:"id"`
				Passes int    `json:"passes"`
				Fails  int    `json:"fails"`
				Name   string `json:"name"`
			} `json:"checks"`
		} `json:"groups"`
		Checks []interface{} `json:"checks"`
	} `json:"root_group"`
	Options struct {
		SummaryTrendStats []string `json:"summaryTrendStats"`
		SummaryTimeUnit   string   `json:"summaryTimeUnit"`
		NoColor           bool     `json:"noColor"`
	} `json:"options"`
	State struct {
		IsStdErrTTY       bool    `json:"isStdErrTTY"`
		TestRunDurationMs float64 `json:"testRunDurationMs"`
		IsStdOutTTY       bool    `json:"isStdOutTTY"`
	} `json:"state"`
	Metrics map[string]interface{} `json:"metrics"`
}

type ZapResult struct {
	Report ZapReport `xml:"report"`
	Html   string    `xml:"html"`
}
type ZapReport struct {
	XMLName   xml.Name  `xml:"OWASPZAPReport"`
	Text      string    `xml:",chardata"`
	Version   string    `xml:"version,attr"`
	Generated string    `xml:"generated,attr"`
	Site      []ZapSite `xml:"site"`
}
type ZapSite struct {
	Text   string `xml:",chardata"`
	Name   string `xml:"name,attr"`
	Host   string `xml:"host,attr"`
	Port   string `xml:"port,attr"`
	Ssl    string `xml:"ssl,attr"`
	Alerts struct {
		Text      string         `xml:",chardata"`
		Alertitem []ZapAlertItem `xml:"alertitem"`
	} `xml:"alerts"`
}
type ZapAlertItem struct {
	Text           string `xml:",chardata"`
	Pluginid       string `xml:"pluginid"`
	AlertRef       string `xml:"alertRef"`
	Alert          string `xml:"alert"`
	Name           string `xml:"name"`
	Riskcode       string `xml:"riskcode"`
	Confidence     string `xml:"confidence"`
	Riskdesc       string `xml:"riskdesc"`
	Confidencedesc string `xml:"confidencedesc"`
	Desc           string `xml:"desc"`
	Instances      struct {
		Text     string `xml:",chardata"`
		Instance []struct {
			Text     string `xml:",chardata"`
			URI      string `xml:"uri"`
			Method   string `xml:"method"`
			Param    string `xml:"param"`
			Attack   string `xml:"attack"`
			Evidence string `xml:"evidence"`
		} `xml:"instance"`
	} `xml:"instances"`
	Count     string `xml:"count"`
	Solution  string `xml:"solution"`
	Otherinfo string `xml:"otherinfo"`
	Reference string `xml:"reference"`
	Cweid     string `xml:"cweid"`
	Wascid    string `xml:"wascid"`
	Sourceid  string `xml:"sourceid"`
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

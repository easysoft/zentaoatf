package commConsts

type ResponseCode struct {
	Code int64  `json:"code"`
	Key  string `json:"message"`
}

var (
	Success  = ResponseCode{0, "request_success"}
	CommErr  = ResponseCode{100, "comm_err"}
	ParamErr = ResponseCode{200, "param_error"}

	NeedInitErr           = ResponseCode{1000, "data_not_init"}
	BizErrWorkspaceConfig = ResponseCode{2000, "workspace_config_err"}
	BizErrNameNotExist    = ResponseCode{3000, "record_not_found_by_name"}
)

type ResultStatus string

const (
	PASS    ResultStatus = "pass"
	FAIL    ResultStatus = "fail"
	SKIP    ResultStatus = "skip"
	BLOCKED ResultStatus = "blocked"
	UNKNOWN ResultStatus = "unknown"
)

func (e ResultStatus) String() string {
	return string(e)
}

type ExecCmd string

const (
	ExecInit ExecCmd = "init"
	ExecStop ExecCmd = "execStop"

	ExecCase   ExecCmd = "execCase"
	ExecModule ExecCmd = "execModule"
	ExecSuite  ExecCmd = "execSuite"
	ExecTask   ExecCmd = "execTask"

	ExecUnit ExecCmd = "execUnit"
)

func (e ExecCmd) String() string {
	return string(e)
}

type WsMsgCategory string

const (
	Communication WsMsgCategory = "communication"
	Exec          WsMsgCategory = "exec"
	Output        WsMsgCategory = "output"
	Error         WsMsgCategory = "error"
	Unknown       WsMsgCategory = ""
)

func (e WsMsgCategory) String() string {
	return string(e)
}

type OsType string

const (
	OsWindows OsType = "windows"
	OsLinux   OsType = "linux"
	OsMac     OsType = "mac"
)

func (e OsType) String() string {
	return string(e)
}
func (OsType) Get(osName string) OsType {
	return OsType(osName)
}

type ExecBy string

const (
	Case   ExecBy = "case"
	Module ExecBy = "module"
	Suite  ExecBy = "suite"
	Task   ExecBy = "task"
)

func (e ExecBy) String() string {
	return string(e)
}
func (ExecBy) Get(str string) ExecBy {
	return ExecBy(str)
}

type TestType string

const (
	TestFunc TestType = "func"
	TestUnit TestType = "unit"
)

func (e TestType) String() string {
	return string(e)
}
func (TestType) Get(str string) TestType {
	return TestType(str)
}

type TestTool string

const (
	ZTF     TestTool = "ztf"
	JUnit   TestTool = "junit"
	TestNG  TestTool = "testng"
	PHPUnit TestTool = "phpunit"
	PyTest  TestTool = "pytest"
	Jest    TestTool = "jest"
	CppUnit TestTool = "cppunit"
	GTest   TestTool = "gtest"
	QTest   TestTool = "qtest"

	AutoIt         TestTool = "autoit"
	Selenium       TestTool = "selenium"
	Appium         TestTool = "appium"
	RobotFramework TestTool = "robotframework"
	Cypress        TestTool = "cypress"
)

func (e TestTool) String() string {
	return string(e)
}
func (TestTool) Get(str string) TestTool {
	return TestTool(str)
}

type BuildTool string

const (
	Maven BuildTool = "maven"
)

func (e BuildTool) String() string {
	return string(e)
}
func (BuildTool) Get(str string) BuildTool {
	return BuildTool(str)
}

type PlatformType string

const (
	Android PlatformType = "android"
	Ios     PlatformType = "ios"
	Host    PlatformType = "host"
	Vm      PlatformType = "vm"
)

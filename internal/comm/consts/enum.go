package commConsts

type ZentaoRequestType string

const (
	PathInfo ZentaoRequestType = "PATH_INFO"
	Get      ZentaoRequestType = "GET"
	Empty    ZentaoRequestType = ""
)

func (e ZentaoRequestType) String() string {
	return string(e)
}

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
)

func (e ExecCmd) String() string {
	return string(e)
}

type WsMsgCategory string

const (
	Communication WsMsgCategory = "communication"
	Exec          WsMsgCategory = "exec"
	Output        WsMsgCategory = "output"
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

type TestType string

const (
	TestFunc TestType = "func"
	TestUnit TestType = "unit"
	TestAuto TestType = "auto"
)

func (e TestType) String() string {
	return string(e)
}
func (TestType) Get(str string) TestType {
	return TestType(str)
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

type UnitTestFramework string

const (
	JUnit   UnitTestFramework = "junit"
	TestNG  UnitTestFramework = "testng"
	PHPUnit UnitTestFramework = "phpunit"
	PyTest  UnitTestFramework = "pytest"
	Jest    UnitTestFramework = "jest"
	CppUnit UnitTestFramework = "cppunit"
	GTest   UnitTestFramework = "gtest"
	QTest   UnitTestFramework = "qtest"
)

func (e UnitTestFramework) String() string {
	return string(e)
}
func (UnitTestFramework) Get(str string) UnitTestFramework {
	return UnitTestFramework(str)
}

type AutoTestTool string

const (
	AutoIt         AutoTestTool = "autoit"
	Selenium       AutoTestTool = "selenium"
	Appium         AutoTestTool = "appium"
	RobotFramework AutoTestTool = "robotframework"
	Cypress        AutoTestTool = "cypress"
)

func (e AutoTestTool) String() string {
	return string(e)
}
func (AutoTestTool) Get(str string) AutoTestTool {
	return AutoTestTool(str)
}

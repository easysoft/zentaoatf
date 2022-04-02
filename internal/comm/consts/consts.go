package commConsts

const (
	App        = "ztf"
	AppServer  = "server"
	AppAgent   = "agent"
	AppCommand = "cmd"

	ConfigVersion = "3.0"
	ConfigDir     = "conf"
	ConfigFile    = "ztf.conf"
	LogDirName    = "log"

	ExtNameSuite = "cs"
	LogText      = "log.txt"
	ResultText   = "result.txt"
	ResultJson   = "result.json"
	ResultZip    = "result.zip"

	ExpectResultPass = "pass"

	PathInfo = "PATH_INFO"
	Get      = "GET"
)

var (
	UnitBuildToolMap = map[string]BuildTool{
		"mvn": Maven,
	}

	AutoTestTypes       = []string{"selenium", "appium"}
	UnitTestTypeJunit   = "junit"
	UnitTestTypeTestNG  = "testng"
	UnitTestTypeRobot   = "robot"
	UnitTestTypeCypress = "cypress"
	UnitTestTypes       = []string{UnitTestTypeJunit, UnitTestTypeTestNG, UnitTestTypeRobot, UnitTestTypeCypress,
		"phpunit", "pytest", "jest", "cppunit", "gtest", "qtest"}
	UnitTestToolMvn   = "mvn"
	UnitTestToolRobot = "robot"
)

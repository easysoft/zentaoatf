package constant

import (
	"fmt"
	"os"
)

const (
	AppName   = "ztf"
	ConfigVer = 2.0
)

var (
	PthSep = string(os.PathSeparator)

	ConfigFile = fmt.Sprintf("conf%s%s.conf", string(os.PathSeparator), AppName)

	UrlZenTaoSettings = "zentaoSettings"
	UrlImportProject  = "importProject"
	UrlSubmitResult   = "submitResults"
	UrlReportBug      = "reportBug"

	ExtNameSuite  = "cs"
	ExtNameJson   = "json"
	ExtNameResult = "txt"

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	EnRes = fmt.Sprintf("res%smessages_en.json", string(os.PathSeparator))
	ZhRes = fmt.Sprintf("res%smessages_zh.json", string(os.PathSeparator))

	LogDir = fmt.Sprintf("log%s", string(os.PathSeparator))

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10

	RequestTypePathInfo = "PATH_INFO"

	AutoTestTypes      = []string{"selenium", "appium"}
	UnitTestTypeJunit  = "junit"
	UnitTestTypeTestNG = "testng"
	UnitTestTypeRobot  = "robot"
	UnitTestTypes      = []string{UnitTestTypeJunit, UnitTestTypeTestNG, UnitTestTypeRobot,
		"phpunit", "pytest", "jest", "cppunit", "gtest", "qtest"}
	UnitTestToolMvn   = "mvn"
	UnitTestToolRobot = "robot"

	RunModeCommon  = "common"
	RunModeServer  = "server"
	RunModeRequest = "request"
)

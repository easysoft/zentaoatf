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

	CmdViewHeight    = 10
	ExpectResultPass = "pass"

	RequestTypePathInfo = "PATH_INFO"

	AutoTestTypes       = []string{"selenium", "appium"}
	UnitTestTypeJunit   = "junit"
	UnitTestTypeTestNG  = "testng"
	UnitTestTypeRobot   = "robot"
	UnitTestTypeCypress = "cypress"
	UnitTestTypes       = []string{UnitTestTypeJunit, UnitTestTypeTestNG, UnitTestTypeRobot, UnitTestTypeCypress,
		"phpunit", "pytest", "jest", "cppunit", "gtest", "qtest"}
	UnitTestToolMvn   = "mvn"
	UnitTestToolRobot = "robot"

	RunModeCommon  = "common"
	RunModeServer  = "server"
	RunModeRequest = "request"

	LangCommentsTagMap = map[string][]string{
		"bat":        {"goto start", ":start"},
		"javascript": {"/\\*{1,}", "\\*{1,}/"},
		"lua":        {"--\\[\\[", "\\]\\]"},
		"perl":       {"=pod", "=cut"},
		"php":        {"/\\*{1,}", "\\*{1,}/"},
		"python":     {"'''", "'''"},
		"ruby":       {"=begin", "=end"},
		"shell":      {":<<!", "!"},
		"tcl":        {"set case {", "}"},
	}

	LangCommentsRegxMap = map[string][]string{
		"bat":        {"^\\s*goto start\\s*$", "^\\s*:start\\s*$"},
		"javascript": {"^\\s*/\\*{1,}\\s*$", "^\\s*\\*{1,}/\\s*$"},
		"lua":        {"^\\s*--\\[\\[\\s*$", "^\\s*\\]\\]\\s*$"},
		"perl":       {"^\\s*=pod\\s*$", "^\\s*=cut\\s*$"},
		"php":        {"^\\s*/\\*{1,}\\s*$", "^\\s*\\*{1,}/\\s*$"},
		"python":     {"^\\s*'''\\s*$", "^\\s*'''\\s*$"},
		"ruby":       {"^\\s*=begin\\s*$", "^\\s*=end\\s*$"},
		"shell":      {"^\\s*:<<!\\s*$", "^\\s*!\\s*$"},
		"tcl":        {"^\\s*set case {", "^\\s*}"},
	}
)

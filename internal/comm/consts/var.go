package commConsts

import (
	"github.com/awesome-gocui/gocui"
)

var (
	LanguageZh = "zh"
	LanguageEn = "en"
	Language   = LanguageZh

	Verbose   = true
	IsRelease bool
	//ExeDir     string
	WorkDir    string
	ExecLogDir string
	LogDir     string

	RequestType string
	ComeFrom    string
)

var (
	Cui              *gocui.Gui
	MainViewHeight   int
	ConfigPath       string
	ServerWorkDir    string
	ServerProjectDir string

	UnitTestType    string
	UnitBuildTool   BuildTool
	UnitTestTool    TestTool
	UnitTestResult  string
	UnitTestResults = "results"
	ProductId       string

	ZenTaoVersion string
	SessionVar    string
	SessionId     string
	RequestFix    string

	ScriptExtToNameMap map[string]string
	CurrScriptFile     string // scripts/tc-001.py
	CurrResultDate     string // 2019-08-15T173802
	CurrCaseId         int    // 2019-08-15T173802

	ScreenWidth    int
	ScreenHeight   int
	CurrBugStepIds string
	Interpreter    string

	// server
	RunMode     string
	IP          string
	MAC         string
	Port        int
	Platform    string
	AgentLogDir string
)

package utils

import (
	"github.com/easysoft/zentaoatf/src/misc"
	"github.com/jroimartin/gocui"
)

const (
	PreferenceFile = "preferences.yaml"
	ConfigFile     = "conf.yaml"

	UrlImportProject = "importProject"
	UrlSubmitResult  = "submitResults"
	UrlReportBug     = "reportBug"

	SuiteExt string = "suite"

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	EnRes = "res/messages_en.json"
	ZhRes = "res/messages_zh.json"

	ScriptDir = "scripts/"
	LogDir    = "logs/"

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10

	CuiRunOutputView = "panelFileContent"
)

var RunMode misc.RunMode
var RunDir string
var RunFromCui bool
var Cui *gocui.Gui
var MainViewHeight int

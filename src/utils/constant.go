package utils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/misc"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/jroimartin/gocui"
	"os"
)

var (
	PreferenceFile = "preferences.yaml"
	ConfigFile     = "conf.yaml"

	UrlZentaoSettings = "zentaoSettings"
	UrlImportProject  = "importProject"
	UrlSubmitResult   = "submitResults"
	UrlReportBug      = "reportBug"

	SuiteExt string = "suite"

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	EnRes = fmt.Sprintf("res%smessages_en.json", string(os.PathSeparator))
	ZhRes = fmt.Sprintf("res%smessages_zh.json", string(os.PathSeparator))

	ScriptDir = fmt.Sprintf("scripts%s", string(os.PathSeparator))
	LogDir    = fmt.Sprintf("logs%s", string(os.PathSeparator))

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10

	CuiRunOutputView = "panelFileContent"

	RequestTypePathInfo = "PATH_INFO"
)

var ZendaoSettings model.ZentaoSettings
var RunMode misc.RunMode
var RunDir string
var RunFromCui bool
var Cui *gocui.Gui
var MainViewHeight int

var SessionVar string
var SessionId string
var RequestType string
var RequestFix string

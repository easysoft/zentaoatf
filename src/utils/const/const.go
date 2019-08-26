package constant

import (
	"fmt"
	"os"
)

var (
	ConfigFile = "conf.yaml"

	UrlZentaoSettings = "zentaoSettings"
	UrlImportProject  = "importProject"
	UrlSubmitResult   = "submitResults"
	UrlReportBug      = "reportBug"

	ExtNameSuite = "suite"
	ExtNameJson  = "json"
	ExtNameTxt   = "txt"

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

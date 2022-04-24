package consts

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"os"
)

const (
	FilePthSep = string(os.PathSeparator)
)

var (
	ConfigFile = fmt.Sprintf("conf%s%s.conf", FilePthSep, commConsts.App)

	ExtNameSuite  = "cs"
	ExtNameJson   = "json"
	ExtNameResult = "txt"

	LogDir = fmt.Sprintf("log%s", FilePthSep)

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight    = 10
	ExpectResultPass = "pass"

	RequestTypePathInfo = "PATH_INFO"

	RunModeCommon  = "common"
	RunModeServer  = "server"
	RunModeRequest = "request"
)

package consts

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"os"
)

const (
	PthSep           = string(os.PathSeparator)
	ConfigFileServer = "server.yaml"
	ConfigFileAgent  = "agent.yaml"
)

var (
	ConfigFile = fmt.Sprintf("conf%s%s.conf", PthSep, commConsts.App)

	ExtNameSuite  = "cs"
	ExtNameJson   = "json"
	ExtNameResult = "txt"

	LogDir = fmt.Sprintf("log%s", PthSep)

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

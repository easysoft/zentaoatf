package domain

import (
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	"time"
)

type Build struct {
	WorkDir    string
	ProjectDir string
	AppPath    string

	ID       uint
	Serial   string
	Priority int
	NodeIp   string
	NodePort int
	DeviceIp string

	BuildType             serverConst.BuildType
	AppiumPort            int
	SeleniumDriverType    serverConst.BrowserType
	SeleniumDriverVersion string
	SeleniumDriverPath    string

	QueueId uint

	ScriptUrl   string
	ScmAddress  string
	ScmAccount  string
	ScmPassword string

	AppUrl          string
	BuildCommands   string
	ResultFiles     string
	KeepResultFiles MyBool
	ResultPath      string
	ResultMsg       string

	StartTime    time.Time
	CompleteTime time.Time

	Progress serverConst.BuildProgress
	Status   serverConst.BuildStatus
}

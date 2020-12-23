package domain

import (
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	"time"
)

type Build struct {
	Debug     bool     `json:"debug,omitempty"`
	ProductId string   `json:"productId,omitempty"`
	SuiteId   string   `json:"suiteId,omitempty"`
	TaskId    string   `json:"taskId,omitempty"`
	Files     []string `json:"files,omitempty"`

	UnitTestType string `json:"unitTestType,omitempty"`
	UnitTestTool string `json:"unitTestTool,omitempty"`
	UnitTestCmd  string `json:"unitTestCmd,omitempty"`

	WorkDir    string `json:"workDir",omitempty`
	ProjectDir string `json:"projectDir,omitempty"`
	AppPath    string `json:"appPath,omitempty"`

	ID           uint   `json:"id,omitempty"`
	QueueId      uint   `json:"queueId,omitempty"`
	Priority     int    `json:"priority,omitempty"`
	NodeIp       string `json:"priority,omitempty"`
	NodePort     int    `json:"nodePort,omitempty"`
	DeviceSerial string `json:"deviceSerial,omitempty"`
	DeviceIp     string `json:"deviceIp,omitempty"`

	BuildType             serverConst.BuildType   `json:"buildType,omitempty"`
	AppiumPort            int                     `json:"appiumPort,omitempty"`
	SeleniumDriverType    serverConst.BrowserType `json:"seleniumDriverType,omitempty"`
	SeleniumDriverVersion string                  `json:"seleniumDriverVersion,omitempty"`
	SeleniumDriverPath    string                  `json:"seleniumDriverPath,omitempty"`

	ScriptUrl   string `json:"scriptUrl,omitempty"`
	ScmAddress  string `json:"scmAddress,omitempty"`
	ScmAccount  string `json:"scmAccount,omitempty"`
	ScmPassword string `json:"scmPassword,omitempty"`

	AppUrl          string `json:"appUrl,omitempty"`
	BuildCommands   string `json:"buildCommands,omitempty"`
	ResultFiles     string `json:"resultFiles,omitempty"`
	KeepResultFiles MyBool `json:"keepResultFiles,omitempty"`
	ResultPath      string `json:"resultPath,omitempty"`
	ResultMsg       string `json:"resultMsg,omitempty"`

	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`

	Progress serverConst.BuildProgress `json:"progress,omitempty"`
	Status   serverConst.BuildStatus   `json:"status,omitempty"`
}

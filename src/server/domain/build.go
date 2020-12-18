package domain

import (
	serverConst "github.com/easysoft/zentaoatf/src/server/utils/const"
	"time"
)

type Build struct {
	Debug     bool     `json:"debug"`
	ProductId string   `json:"productId"`
	SuiteId   string   `json:"suiteId"`
	TaskId    string   `json:"taskId"`
	Files     []string `json:"files"`

	UnitTestType string `json:"unitTestType"`
	UnitTestTool string `json:"unitTestTool"`
	UnitTestCmd  string `json:"unitTestCmd"`

	WorkDir    string `json:"workDir"`
	ProjectDir string `json:"projectDir"`
	AppPath    string `json:"appPath"`

	ID           uint   `json:"id"`
	QueueId      uint   `json:"queueId"`
	Priority     int    `json:"priority"`
	NodeIp       string `json:"priority"`
	NodePort     int    `json:"nodePort"`
	DeviceSerial string `json:"deviceSerial"`
	DeviceIp     string `json:"deviceIp"`

	BuildType             serverConst.BuildType   `json:"buildType"`
	AppiumPort            int                     `json:"appiumPort"`
	SeleniumDriverType    serverConst.BrowserType `json:"seleniumDriverType"`
	SeleniumDriverVersion string                  `json:"seleniumDriverVersion"`
	SeleniumDriverPath    string                  `json:"seleniumDriverPath"`

	ScriptUrl   string `json:"scriptUrl"`
	ScmAddress  string `json:"scmAddress"`
	ScmAccount  string `json:"scmAccount"`
	ScmPassword string `json:"scmPassword"`

	AppUrl          string `json:"appUrl"`
	BuildCommands   string `json:"buildCommands"`
	ResultFiles     string `json:"resultFiles"`
	KeepResultFiles MyBool `json:"keepResultFiles"`
	ResultPath      string `json:"resultPath"`
	ResultMsg       string `json:"resultMsg"`

	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`

	Progress serverConst.BuildProgress `json:"progress"`
	Status   serverConst.BuildStatus   `json:"status"`
}

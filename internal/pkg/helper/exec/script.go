package execHelper

import (
	"fmt"
	"sync"
	"time"

	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

func ExecScript(scriptFile, workspacePath string,
	conf commDomain.WorkspaceConf,
	report *commDomain.ZtfReport, scriptIdx,
	total, pathMaxWidth, numbMaxWidth, titleMaxWidth int,
	ch chan int, wsMsg *websocket.Message, lock *sync.Mutex) {

	key := stringUtils.Md5(scriptFile)

	startTime := time.Now()

	startMsg := i118Utils.Sprintf("start_execution", scriptFile, dateUtils.DateTimeStr(startTime))

	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg(startMsg, "", commConsts.Run,
			iris.Map{"key": key, "status": "start"}, wsMsg)

		logUtils.ExecConsolef(-1, startMsg)
	}
	logUtils.ExecFilef(startMsg)

	logs := ""
	stdOutput, errOutput := RunFile(scriptFile, workspacePath, conf, ch, wsMsg, scriptIdx)
	// stdOutput = strings.Trim(stdOutput, "\n")

	if stdOutput != "" {
		logs = stdOutput
	}
	if errOutput != "" {
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendOutputMsg(errOutput, "", iris.Map{"key": key}, wsMsg)
		}
		logUtils.ExecConsolef(-1, errOutput)
		logUtils.ExecFilef(errOutput)
	}

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))
	report.WorkspacePath = workspacePath

	CheckCaseResult(scriptFile, logs, report, total, secs, pathMaxWidth, numbMaxWidth, titleMaxWidth, wsMsg, errOutput, lock)

	endMsg := i118Utils.Sprintf("end_execution", scriptFile, dateUtils.DateTimeStr(entTime))
	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg(endMsg, "", commConsts.Run, iris.Map{"key": key, "status": "end"}, wsMsg)

		logUtils.ExecConsolef(-1, endMsg)
	}
	logUtils.ExecFilef(endMsg)
}

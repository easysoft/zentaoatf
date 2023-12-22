package execHelper

import (
	"fmt"
	"strings"
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

func ExecScript(execParams commDomain.ExecParams, ch chan int, wsMsg *websocket.Message, lock *sync.Mutex) {

	key := stringUtils.Md5(execParams.ScriptFile)

	startTime := time.Now()

	startMsg := i118Utils.Sprintf("start_execution", execParams.ScriptFile, dateUtils.DateTimeStr(startTime))

	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg(startMsg, "", commConsts.Run,
			iris.Map{"key": key, "status": "start"}, wsMsg)

		logUtils.ExecConsolef(-1, startMsg)
	}
	logUtils.ExecFilef(startMsg)

	logs := ""
	stdOutput, errOutput := RunFile(execParams.ScriptFile, execParams.WorkspacePath, execParams.Conf, ch, wsMsg, execParams.ScriptIdx)
	stdOutput = strings.TrimLeft(stdOutput, "\n")

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
	execParams.Secs = fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))
	execParams.Report.WorkspacePath = execParams.WorkspacePath

	CheckCaseResult(execParams, logs, wsMsg, errOutput, lock)

	endMsg := i118Utils.Sprintf("end_execution", execParams.ScriptFile, dateUtils.DateTimeStr(entTime))
	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg(endMsg, "", commConsts.Run, iris.Map{"key": key, "status": "end"}, wsMsg)

		logUtils.ExecConsolef(-1, endMsg)
	}
	logUtils.ExecFilef(endMsg)
}

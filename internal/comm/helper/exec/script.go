package execHelper

import (
	"fmt"
	"strings"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	websocketHelper "github.com/easysoft/zentaoatf/internal/comm/helper/websocket"
	dateUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/date"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/string"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

func ExecScript(scriptFile, workspacePath string,
	conf commDomain.WorkspaceConf,
	report *commDomain.ZtfReport, scriptIdx,
	total, pathMaxWidth, numbMaxWidth int,
	ch chan int, wsMsg *websocket.Message) {

	key := stringUtils.Md5(scriptFile)

	startTime := time.Now()

	startMsg := i118Utils.Sprintf("start_execution", scriptFile, dateUtils.DateTimeStr(startTime))

	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(startMsg, "", commConsts.Run,
			iris.Map{"key": key, "status": "start"}, wsMsg)

		logUtils.ExecConsolef(-1, startMsg)
	}
	logUtils.ExecFilef(startMsg)

	logs := ""
	stdOutput, errOutput := RunFile(scriptFile, workspacePath, conf, ch, wsMsg)
	stdOutput = strings.Trim(stdOutput, "\n")

	if stdOutput != "" {
		logs = stdOutput
	}
	if errOutput != "" {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(errOutput, "", iris.Map{"key": key}, wsMsg)
		}
		logUtils.ExecConsolef(-1, errOutput)
		logUtils.ExecFilef(errOutput)
	}

	entTime := time.Now()
	secs := fmt.Sprintf("%.2f", float32(entTime.Sub(startTime)/time.Second))
	report.WorkspacePath = workspacePath
	CheckCaseResult(scriptFile, logs, report, scriptIdx, total, secs, pathMaxWidth, numbMaxWidth, wsMsg)

	endMsg := i118Utils.Sprintf("end_execution", scriptFile, dateUtils.DateTimeStr(entTime))
	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(endMsg, "", commConsts.Run, iris.Map{"key": key, "status": "end"}, wsMsg)

		logUtils.ExecConsolef(-1, endMsg)
	}
	logUtils.ExecFilef(endMsg)

	//for i := 0; i < 100000; i++ {
	//	websocketHelper.SendExecMsg(fmt.Sprintf("------%d", i), "", commConsts.Result,
	//		iris.Map{"key": key, "status": "end"}, wsMsg)
	//	time.Sleep(time.Millisecond * 100)
	//}
}

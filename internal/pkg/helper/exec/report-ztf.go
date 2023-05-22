package execHelper

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
)

func GenZTFTestReport(report commDomain.ZtfReport, pathMaxWidth int,
	workspacePath string, ch chan int, wsMsg *websocket.Message) {
	SetRunningIfNeed(ch)

	// print failed case
	failedCount, failedCaseLinesWithCheckpoint := GenZTFReport(report)

	logUtils.ExecConsolef(-1, "")

	if failedCount > 0 {
		msgFail := genZTFFailedMsg(failedCaseLinesWithCheckpoint)

		logUtils.ExecConsolef(-1, msgFail)
		logUtils.ExecFile(msgFail)
	}

	// 生成统计行
	runResult, msgRunColor := GenRunResult(report)

	// 执行%d个用例，耗时%d秒%s。%s，%s，%s。
	// Run %d script in %d sec, %s, %s, %s.
	msgRun := dateUtils.DateTimeStr(time.Now()) + " " + runResult

	websocketHelper.SendExecMsgIfNeed(msgRunColor, "", commConsts.Run, nil, wsMsg)

	logUtils.ExecResult(msgRun)

	resultPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultText)
	msgReport := i118Utils.Sprintf("run_report") + " " + resultPath + "."
	if commConsts.ExecFrom == commConsts.FromCmd {
		msgReport = color.New(color.Bold, color.FgHiWhite).Sprint(i118Utils.Sprintf("run_report")) + " [" + resultPath + "]"
	}

	logUtils.ExecConsole(-1, msgReport)
	logUtils.ExecConsole(-1, msgRun)
	logUtils.ExecResult(msgReport)
	websocketHelper.SendExecMsgIfNeed(msgReport, "", commConsts.Run, map[string]interface{}{
		"logDir": commConsts.ExecLogDir,
	}, wsMsg)
	report.Log = fileUtils.ReadFile(filepath.Join(commConsts.ExecLogDir, commConsts.LogText))

	//report.ProductId, _ = strconv.Atoi(vari.ProductId)
	json, _ := json.MarshalIndent(report, "", "\t")
	jsonPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultJson)
	fileUtils.WriteFile(jsonPath, string(json))
}

func genZTFFailedMsg(failedCaseLinesWithCheckpoint []string) (msgFail string) {
	divider := shellUtils.GenFullScreenDivider()

	msgFail = divider
	msgFail += "\n" + color.New(color.Bold, color.FgHiWhite).Sprint(i118Utils.Sprintf("failed_scripts")) + "\n"
	msgFail += strings.Join(failedCaseLinesWithCheckpoint, "\n")
	msgFail += "\n\n" + divider

	return msgFail
}

func appendFailedStepResult(cs commDomain.FuncResult, failedSteps *[]string) (passStepCount, failedCount int) {
	if len(cs.Steps) == 0 {
		*failedSteps = append(*failedSteps, "   "+i118Utils.Sprintf("no_checkpoints"))
		return
	}

	for _, step := range cs.Steps {
		if step.Status == commConsts.PASS {
			passStepCount++
			continue
		}

		step.Id = strings.TrimRight(step.Id, ".")
		status := GenStatusTxt(cs.Status)

		*failedSteps = append(*failedSteps, fmt.Sprintf("%s %s [%s]", i118Utils.Sprintf("step_prefix", step.Id), status, step.Name))

		for idx1, cp := range step.CheckPoints {
			//cpStatus := commonUtils.BoolToPass(step.Status)
			*failedSteps = append(*failedSteps, fmt.Sprintf("[%s]  %s", i118Utils.Sprintf("expect"), cp.Expect))
			*failedSteps = append(*failedSteps, fmt.Sprintf("[%s]  %s", i118Utils.Sprintf("actual"), cp.Actual))

			if idx1 < len(step.CheckPoints)-1 {
				*failedSteps = append(*failedSteps, "")
			}
		}
		failedCount++
	}

	return
}

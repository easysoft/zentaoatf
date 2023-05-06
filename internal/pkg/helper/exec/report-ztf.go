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
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
)

func GenZTFTestReport(report commDomain.ZtfReport, pathMaxWidth int,
	workspacePath string, ch chan int, wsMsg *websocket.Message) {
	select {
	case _, ok := <-ch:
		if !ok {
			SetRunning(false)
			return
		}
	default:
	}
	// print failed case
	failedCount := 0
	failedCaseLines := make([]string, 0)
	failedCaseLinesWithCheckpoint := make([]string, 0)

	for _, csResult := range report.FuncResult {
		if report.ProductId == 0 && csResult.ProductId > 0 {
			report.ProductId = csResult.ProductId
		}

		if csResult.Status == "fail" {
			if failedCount > 0 {
				failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
			}
			failedCount++

			path := csResult.Path
			lent := runewidth.StringWidth(path)

			relativePath := csResult.Path
			if strings.Contains(relativePath, "/module/") {
				relativePath = relativePath[strings.Index(relativePath, "/module/"):]
			}

			if pathMaxWidth > lent {
				postFix := strings.Repeat(" ", pathMaxWidth-lent)
				path += postFix
				relativePath += postFix
			}

			line := fmt.Sprintf("[%s] %d.%s", relativePath, csResult.Id, csResult.Title)
			trimLine := fmt.Sprintf("[%s] %d.%s", strings.TrimSpace(relativePath), csResult.Id, csResult.Title)
			failedCaseLines = append(failedCaseLines, line)
			failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, trimLine)

			appendFailedStepResult(csResult, &failedCaseLinesWithCheckpoint)
		}
	}
	if failedCount > 0 {
		msgFail := "\n" + i118Utils.Sprintf("failed_scripts") + "\n"
		msgFail += strings.Join(failedCaseLines, "\n") + "\n\n"
		msgFail += strings.Join(failedCaseLinesWithCheckpoint, "\n")

		logUtils.ExecConsolef(color.FgRed, msgFail)
		logUtils.ExecFile(msgFail)
	}

	// 生成统计行
	fmtStr := "%d(%.1f%%) %s"
	passRate := 0
	failRate := 0
	skipRate := 0
	if report.Total > 0 {
		passRate = report.Pass * 100 / report.Total
		failRate = report.Fail * 100 / report.Total
		skipRate = report.Skip * 100 / report.Total
	}

	passStr := fmt.Sprintf(fmtStr, report.Pass, float32(passRate), i118Utils.Sprintf("pass"))
	failStr := fmt.Sprintf(fmtStr, report.Fail, float32(failRate), i118Utils.Sprintf("fail"))
	skipStr := fmt.Sprintf(fmtStr, report.Skip, float32(skipRate), i118Utils.Sprintf("skip"))

	// 执行%d个用例，耗时%d秒%s。%s，%s，%s。
	// Run %d script in %d sec, %s, %s, %s.
	msgRun := dateUtils.DateTimeStr(time.Now()) + " " +
		i118Utils.Sprintf("run_result", report.Total, report.Duration, passStr, failStr, skipStr)

	if commConsts.ExecFrom == commConsts.FromClient {
		msgRunColor := i118Utils.Sprintf("run_result", report.Total, report.Duration,
			fmt.Sprintf(`<span class="result-pass">%s</span>`, passStr),
			fmt.Sprintf(`<span class="result-fail">%s</span>`, failStr),
			fmt.Sprintf(`<span class="result-skip">%s</span>`, skipStr),
		)
		websocketHelper.SendExecMsg(msgRunColor, "", commConsts.Run, nil, wsMsg)
	}

	logUtils.ExecConsole(color.FgCyan, msgRun)
	logUtils.ExecResult(msgRun)

	resultPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultText)
	msgReport := "                    " + i118Utils.Sprintf("run_report", resultPath) + "\n"

	logUtils.ExecConsole(color.FgCyan, msgReport)
	logUtils.ExecResult(msgReport)
	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg(msgReport, "", commConsts.Run, map[string]interface{}{
			"logDir": commConsts.ExecLogDir,
		}, wsMsg)
	}
	report.Log = fileUtils.ReadFile(filepath.Join(commConsts.ExecLogDir, commConsts.LogText))
	report.TestCommand = strings.Join(commConsts.CommandArgs, " ")

	//report.ProductId, _ = strconv.Atoi(vari.ProductId)
	json, _ := json.MarshalIndent(report, "", "\t")
	jsonPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultJson)
	fileUtils.WriteFile(jsonPath, string(json))
}

func appendFailedStepResult(cs commDomain.FuncResult, failedSteps *[]string) (passStepCount, failedCount int) {
	if len(cs.Steps) > 0 {
		for _, step := range cs.Steps {
			if step.Status == commConsts.PASS {
				passStepCount++
				continue
			}

			step.Id = strings.TrimRight(step.Id, ".")
			status := i118Utils.Sprintf(string(step.Status))
			*failedSteps = append(*failedSteps, fmt.Sprintf("  Step %s: %s [%s]", step.Id, status, step.Name))

			for idx1, cp := range step.CheckPoints {
				//cpStatus := commonUtils.BoolToPass(step.Status)
				*failedSteps = append(*failedSteps, fmt.Sprintf("    [Expect] %s", cp.Expect))
				*failedSteps = append(*failedSteps, fmt.Sprintf("    [Actual] %s", cp.Actual))

				if idx1 < len(step.CheckPoints)-1 {
					*failedSteps = append(*failedSteps, "")
				}
			}
			failedCount++
		}
	} else {
		*failedSteps = append(*failedSteps, "   "+i118Utils.Sprintf("no_checkpoints"))
	}

	return
}

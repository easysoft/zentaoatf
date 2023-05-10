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
	failedCaseLinesWithCheckpoint := make([]string, 0)

	for index, csResult := range report.FuncResult {
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

			if pathMaxWidth > lent {
				postFix := strings.Repeat(" ", pathMaxWidth-lent)
				path += postFix
			}

			prefix := i118Utils.Sprintf("test_case_prefix", index+1)
			line := fmt.Sprintf("%s[%s] [%d.%s]", prefix, csResult.Path, csResult.Id, csResult.Title)
			failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, line)

			appendFailedStepResult(csResult, &failedCaseLinesWithCheckpoint)
		}
	}

	logUtils.ExecConsolef(-1, "")

	if failedCount > 0 {
		divider := "--------------------------------"
		window := shellUtils.WindowSize()
		if window.Col != 0 {
			divider = strings.Repeat("-", int(window.Col))
		}

		msgFail := divider
		msgFail += "\n" + color.New(color.Bold, color.FgWhite).Sprint(i118Utils.Sprintf("failed_scripts")) + "\n"
		msgFail += strings.Join(failedCaseLinesWithCheckpoint, "\n")
		msgFail += "\n\n" + divider

		logUtils.ExecConsolef(-1, msgFail)
		logUtils.ExecFile(msgFail)
	}

	// 生成统计行
	fmtStr := "%s%d(%.1f%%)"
	passRate := 0
	failRate := 0
	skipRate := 0
	if report.Total > 0 {
		passRate = report.Pass * 100 / report.Total
		failRate = report.Fail * 100 / report.Total
		skipRate = report.Skip * 100 / report.Total
	}

	passStr := fmt.Sprintf(fmtStr, i118Utils.Sprintf("pass_num"), report.Pass, float32(passRate))
	failStr := fmt.Sprintf(fmtStr, i118Utils.Sprintf("fail_num"), report.Fail, float32(failRate))
	skipStr := fmt.Sprintf(fmtStr, i118Utils.Sprintf("skip_num"), report.Skip, float32(skipRate))

	if commConsts.ExecFrom == commConsts.FromCmd {
		passStr = fmt.Sprintf(fmtStr, color.New(color.FgGreen).Sprint(i118Utils.Sprintf("pass_num")), report.Pass, float32(passRate))
		failStr = fmt.Sprintf(fmtStr, color.New(color.FgRed).Sprint(i118Utils.Sprintf("fail_num")), report.Fail, float32(failRate))
		skipStr = fmt.Sprintf(fmtStr, color.New(color.FgYellow).Sprint(i118Utils.Sprintf("skip_num")), report.Skip, float32(skipRate))
	}

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

	logUtils.ExecResult(msgRun)

	resultPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultText)
	msgReport := i118Utils.Sprintf("run_report", resultPath)

	logUtils.ExecConsole(-1, msgReport)
	logUtils.ExecConsole(-1, msgRun)
	logUtils.ExecResult(msgReport)
	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg(msgReport, "", commConsts.Run, map[string]interface{}{
			"logDir": commConsts.ExecLogDir,
		}, wsMsg)
	}
	report.Log = fileUtils.ReadFile(filepath.Join(commConsts.ExecLogDir, commConsts.LogText))

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
			if commConsts.ExecFrom == commConsts.FromCmd {
				if step.Status == commConsts.FAIL {
					status = color.New(color.FgRed).Sprint(status)
				} else if step.Status == commConsts.PASS {
					status = color.New(color.FgGreen).Sprint(status)
				} else {
					status = color.New(color.FgYellow).Sprint(status)
				}
			}

			*failedSteps = append(*failedSteps, fmt.Sprintf("%s%s: %s [%s]", i118Utils.Sprintf("step"), step.Id, status, step.Name))

			for idx1, cp := range step.CheckPoints {
				//cpStatus := commonUtils.BoolToPass(step.Status)
				*failedSteps = append(*failedSteps, fmt.Sprintf("[%s] %s", i118Utils.Sprintf("expect"), cp.Expect))
				*failedSteps = append(*failedSteps, fmt.Sprintf("[%s] %s", i118Utils.Sprintf("actual"), cp.Actual))

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

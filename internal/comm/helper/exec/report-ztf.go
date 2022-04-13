package execHelper

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	websocketHelper "github.com/easysoft/zentaoatf/internal/comm/helper/websocket"
	dateUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/date"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
	"path/filepath"
	"strings"
	"time"
)

func GenZTFTestReport(report commDomain.ZtfReport, pathMaxWidth int,
	workspacePath string, wsMsg *websocket.Message) {

	// print failed case
	failedCount := 0
	failedCaseLines := make([]string, 0)
	failedCaseLinesWithCheckpoint := make([]string, 0)

	for _, cs := range report.FuncResult {
		if report.ProductId == 0 && cs.ProductId > 0 {
			report.ProductId = cs.ProductId
		}

		if cs.Status == "fail" {
			if failedCount > 0 {
				failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
			}
			failedCount++

			path := cs.Path
			lent := runewidth.StringWidth(path)

			if pathMaxWidth > lent {
				postFix := strings.Repeat(" ", pathMaxWidth-lent)
				path += postFix
			}

			line := fmt.Sprintf("[%s] %d.%s", cs.Path, cs.Id, cs.Title)
			failedCaseLines = append(failedCaseLines, line)
			failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, line)

			if len(cs.Steps) > 0 {
				stepNumb := 0
				for _, step := range cs.Steps {
					if step.Status == commConsts.PASS {
						continue
					}

					if stepNumb > 0 {
						failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
					}
					stepNumb++

					step.Id = strings.TrimRight(step.Id, ".")
					status := i118Utils.Sprintf(string(step.Status))
					failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("Step %s: %s", step.Id, status))

					for idx1, cp := range step.CheckPoints {
						//cpStatus := commonUtils.BoolToPass(step.Status)
						failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("[Expect] %s", cp.Expect))
						failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("[Actual] %s", cp.Actual))

						if idx1 < len(step.CheckPoints)-1 {
							failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
						}
					}
				}
			} else {
				failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "   "+i118Utils.Sprintf("no_checkpoints"))
			}
		}
	}
	if failedCount > 0 {
		msgFail := "\n" + i118Utils.Sprintf("failed_scripts")
		msgFail += strings.Join(failedCaseLines, "\n")
		msgFail += strings.Join(failedCaseLinesWithCheckpoint, "\n")

		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendExecMsg(msgFail, "", commConsts.Error, wsMsg)
		}

		logUtils.ExecConsolef(color.FgRed, msgFail)
		logUtils.ExecFile(msgFail)
	}

	// 生成统计行
	secTag := ""
	if commConsts.Language == "en" && report.Duration > 1 {
		secTag = "s"
	}

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

	// 执行%d个用例，耗时%d秒%s。%s，%s，%s。报告%s。
	msgRun := dateUtils.DateTimeStr(time.Now()) + " " +
		i118Utils.Sprintf("run_result",
			report.Total, report.Duration, secTag,
			passStr, failStr, skipStr,
		)

	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(msgRun, "", commConsts.Run, wsMsg)
	}

	logUtils.ExecConsole(color.FgCyan, msgRun)
	logUtils.ExecResult(msgRun)

	resultPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultText)
	msgReport := "                    " + i118Utils.Sprintf("run_report", resultPath) + "\n"

	if commConsts.ExecFrom != commConsts.FromCmd {
		websocketHelper.SendExecMsg(msgReport, "false", commConsts.Run, wsMsg)
	}

	logUtils.ExecConsole(color.FgCyan, msgReport)
	logUtils.ExecResult(msgReport)

	//report.ProductId, _ = strconv.Atoi(vari.ProductId)
	json, _ := json.MarshalIndent(report, "", "\t")
	jsonPath := filepath.Join(commConsts.ExecLogDir, commConsts.ResultJson)
	fileUtils.WriteFile(jsonPath, string(json))
}

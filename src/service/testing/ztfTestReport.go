package testingService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"strings"
	"time"
)

func GenZtfTestReport(report model.TestReport, pathMaxWidth int) {
	if len(report.Cases) == 0 {
		return
	}

	// print failed case
	failedCount := 0
	failedCaseLines := make([]string, 0)
	failedCaseLinesWithCheckpoint := make([]string, 0)

	for _, cs := range report.Cases {
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
					if step.Status {
						continue
					}

					if stepNumb > 0 {
						failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
					}
					stepNumb++

					step.Id = strings.TrimRight(step.Id, ".")
					status := i118Utils.I118Prt.Sprintf(commonUtils.BoolToPass(step.Status))
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
				failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "   "+i118Utils.I118Prt.Sprintf("no_checkpoints"))
			}
		}
	}
	if failedCount > 0 {
		logUtils.ScreenAndResult("\n" + i118Utils.I118Prt.Sprintf("failed_scripts"))
		logUtils.Screen(strings.Join(failedCaseLines, "\n"))
		logUtils.Result(strings.Join(failedCaseLinesWithCheckpoint, "\n"))
	}

	secTag := ""
	if vari.Config.Language == "en" && report.Duration > 1 {
		secTag = "s"
	}

	// 生成统计行
	fmtStr := "%d(%.1f%%) %s"
	passStr := fmt.Sprintf(fmtStr, report.Pass, float32(report.Pass*100/report.Total), i118Utils.I118Prt.Sprintf("pass"))
	failStr := fmt.Sprintf(fmtStr, report.Fail, float32(report.Fail*100/report.Total), i118Utils.I118Prt.Sprintf("fail"))
	skipStr := fmt.Sprintf(fmtStr, report.Skip, float32(report.Skip*100/report.Total), i118Utils.I118Prt.Sprintf("skip"))

	// 带映带结果文件
	logUtils.Result("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("run_scripts",
			report.Total, report.Duration, secTag,
			passStr, failStr, skipStr,
			vari.LogDir+"result.txt",
		))
	// 打印到屏幕
	logUtils.Screen("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("run_scripts",
			report.Total, report.Duration, secTag,
			color.GreenString(passStr), color.RedString(failStr), color.YellowString(skipStr),
			vari.LogDir+"result.txt",
		))

	//println("===" + vari.LogDir)

	json, _ := json.Marshal(report)
	fileUtils.WriteFile(vari.LogDir+"result.json", string(json))
}

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
	"github.com/mattn/go-runewidth"
	"strings"
	"time"
)

func Report(report model.TestReport, pathMaxWidth int) {
	// print failed case
	failedCount := 0
	failedCaseLines := make([]string, 0)
	failedCaseLinesWithCheckpoint := make([]string, 0)
	for _, cs := range report.Cases {
		if cs.Status == "fail" {
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
		}

		if len(cs.Steps) > 0 {
			for idx, step := range cs.Steps {
				status := i118Utils.I118Prt.Sprintf(commonUtils.BoolToPass(step.Status))
				failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("  %s %s", step.Id, status))

				for idx1, cp := range step.CheckPoints {
					//cpStatus := commonUtils.BoolToPass(step.Status)
					failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("    [Expect] %s", cp.Expect))
					failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("    [Actual] %s", cp.Actual))

					if idx1 < len(step.CheckPoints)-1 {
						failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
					}
				}

				if idx < len(cs.Steps)-1 {
					failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
				}
			}
		} else {
			failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "   "+i118Utils.I118Prt.Sprintf("no_checkpoints"))
		}
	}
	if failedCount > 0 {
		logUtils.ScreenAndResult("\n" + i118Utils.I118Prt.Sprintf("failed_scripts"))
		logUtils.Screen(strings.Join(failedCaseLines, "\n"))
		logUtils.Result(strings.Join(failedCaseLinesWithCheckpoint, "\n"))
	}

	logUtils.ScreenAndResult("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("run_scripts",
			report.Total, report.Duration,
			report.Pass, float32(report.Pass*100/report.Total), i118Utils.I118Prt.Sprintf("pass"),
			report.Fail, float32(report.Fail*100/report.Total), i118Utils.I118Prt.Sprintf("fail"),
			report.Skip, float32(report.Skip*100/report.Total), i118Utils.I118Prt.Sprintf("skip"),
			vari.ZtfDir+"result.txt",
		))

	json, _ := json.Marshal(report)
	fileUtils.WriteFile(vari.LogDir+"result.json", string(json))
}

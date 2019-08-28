package testingService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"strings"
	"time"
)

func Print(report model.TestReport) {
	startSec := time.Unix(report.StartTime, 0)
	endSec := time.Unix(report.EndTime, 0)

	logs := make([]string, 0)

	logUtils.PrintAndLog(&logs, i118Utils.I118Prt.Sprintf("run_scripts", vari.WorkDir, report.Env))

	logUtils.PrintAndLog(&logs, i118Utils.I118Prt.Sprintf("time_from_to",
		startSec.Format("2006-01-02 15:04:05"), endSec.Format("2006-01-02 15:04:05"), report.Duration))

	logUtils.PrintAndLog(&logs, fmt.Sprintf("%s: %d", i118Utils.I118Prt.Sprintf("total"), report.Total))
	logUtils.PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", i118Utils.I118Prt.Sprintf("pass"), report.Pass), color.FgGreen)
	logUtils.PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", i118Utils.I118Prt.Sprintf("fail"), report.Fail), color.FgRed)
	logUtils.PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", i118Utils.I118Prt.Sprintf("skip"), report.Skip), color.FgYellow)

	for _, cs := range report.Cases {
		str := "\n %s %d: %s"
		status := cs.Status
		statusColor := logUtils.ColoredStatus(status)

		logs = append(logs, fmt.Sprintf(str, status, cs.Id, cs.Path))
		logUtils.Printt(fmt.Sprintf(str+"\n", statusColor, cs.Id, cs.Path))

		if len(cs.Steps) > 0 {
			count := 0
			for _, step := range cs.Steps {
				if count > 0 { // 空行
					logUtils.PrintAndLog(&logs, "")
				}

				str := "  %s %d:   %s"
				status := commonUtils.BoolToPass(step.Status)
				statusColor := logUtils.ColoredStatus(status)

				logs = append(logs, fmt.Sprintf(str, status, step.Id, step.Name))
				logUtils.Printt(fmt.Sprintf(str, statusColor, step.Id, step.Name+"\n"))

				count1 := 0
				for _, cp := range step.CheckPoints {
					if count1 > 0 { // 空行
						logUtils.PrintAndLog(&logs, "")
					}

					cpStatus := commonUtils.BoolToPass(step.Status)
					cpStatusColored := logUtils.ColoredStatus(cpStatus)
					logs = append(logs, fmt.Sprintf("    %s %d: %s", commonUtils.BoolToPass(cp.Status), cp.Numb,
						i118Utils.I118Prt.Sprintf("checkpoint")))
					logUtils.Printt(fmt.Sprintf("    %s %d: %s\n", cpStatusColored, cp.Numb, i118Utils.I118Prt.Sprintf("checkpoint")))

					logUtils.PrintAndLog(&logs, fmt.Sprintf("      %s %s", i118Utils.I118Prt.Sprintf("expect_result"), cp.Expect))
					logUtils.PrintAndLog(&logs, fmt.Sprintf("      %s %s", i118Utils.I118Prt.Sprintf("actual_result"), cp.Actual))

					count1++
				}

				count++
			}
		} else {
			logUtils.PrintAndLog(&logs, "   "+i118Utils.I118Prt.Sprintf("no_checkpoints"))
		}
	}

	fileUtils.WriteFile(constant.LogDir+vari.RunDir+"result.txt", strings.Join(logs, "\n"))

	json, _ := json.Marshal(report)
	fileUtils.WriteFile(constant.LogDir+vari.RunDir+"result.json", string(json))
}

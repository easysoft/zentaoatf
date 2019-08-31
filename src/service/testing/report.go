package testingService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"time"
)

func Print(report model.TestReport) {
	startSec := time.Unix(report.StartTime, 0)
	endSec := time.Unix(report.EndTime, 0)

	logUtils.PrintAndLog(i118Utils.I118Prt.Sprintf("run_scripts", vari.WorkDir, report.Env))

	logUtils.PrintAndLog(i118Utils.I118Prt.Sprintf("time_from_to",
		startSec.Format("2006-01-02 15:04:05"), endSec.Format("2006-01-02 15:04:05"), report.Duration))

	logUtils.PrintAndLog(fmt.Sprintf("%s: %d", i118Utils.I118Prt.Sprintf("total"), report.Total))
	logUtils.PrintAndLogColorLn(fmt.Sprintf("  %s: %d", i118Utils.I118Prt.Sprintf("pass"), report.Pass), color.FgGreen)
	logUtils.PrintAndLogColorLn(fmt.Sprintf("  %s: %d", i118Utils.I118Prt.Sprintf("fail"), report.Fail), color.FgRed)
	logUtils.PrintAndLogColorLn(fmt.Sprintf("  %s: %d", i118Utils.I118Prt.Sprintf("skip"), report.Skip), color.FgYellow)

	for _, cs := range report.Cases {
		str := "\n%d. %s %s, %s"
		status := cs.Status
		statusColor := logUtils.ColoredStatus(status)

		//logs = append(logs, fmt.Sprintf(str, status, cs.Id, cs.Path))
		logUtils.PrintTo(fmt.Sprintf(str+"\n", cs.Id, cs.Title, statusColor, cs.Path))

		if len(cs.Steps) > 0 {
			count := 0
			for _, step := range cs.Steps {
				if count > 0 { // 空行
					logUtils.PrintAndLog("")
				}

				str := "[Step%d]: %s \n"
				//status := commonUtils.BoolToPass(step.Status)
				//statusColor := logUtils.ColoredStatus(status)

				//logs = append(logs, fmt.Sprintf(str, status, step.Name))
				logUtils.PrintTo(fmt.Sprintf(str, step.Id, step.Name))

				count1 := 0
				for _, cp := range step.CheckPoints {
					if count1 > 0 { // 空行
						logUtils.PrintAndLog("")
					}

					//cpStatus := commonUtils.BoolToPass(step.Status)
					//cpStatusColored := logUtils.ColoredStatus(cpStatus)
					//logs = append(logs, fmt.Sprintf("    %s : %s %d", commonUtils.BoolToPass(cp.Status),
					//	i118Utils.I118Prt.Sprintf("checkpoint"), cp.Numb))
					//logUtils.PrintTo(fmt.Sprintf("    %s: %s %d\n", cpStatusColored, i118Utils.I118Prt.Sprintf("checkpoint"), cp.Numb))

					logUtils.PrintAndLog(fmt.Sprintf("[Expect] %s", cp.Expect))
					logUtils.PrintAndLog(fmt.Sprintf("[Actual] %s", cp.Actual))

					count1++
				}

				count++
			}
		} else {
			logUtils.PrintAndLog("   " + i118Utils.I118Prt.Sprintf("no_checkpoints"))
		}
	}

	json, _ := json.Marshal(report)
	fileUtils.WriteFile(constant.LogDir+vari.RunDir+"result.json", string(json))
}

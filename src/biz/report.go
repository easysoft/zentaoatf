package biz

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"strings"
	"time"
)

func Print(report model.TestReport, workDir string) {
	startSec := time.Unix(report.StartTime, 0)
	endSec := time.Unix(report.EndTime, 0)

	logs := make([]string, 0)

	PrintAndLog(&logs, utils.I118Prt.Sprintf("run_scripts", report.Path, report.Env))

	PrintAndLog(&logs, utils.I118Prt.Sprintf("time_from_to",
		startSec.Format("2006-01-02 15:04:05"), endSec.Format("2006-01-02 15:04:05"), report.Duration))

	PrintAndLog(&logs, fmt.Sprintf("%s: %d", utils.I118Prt.Sprintf("total"), report.Total))
	PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", utils.I118Prt.Sprintf("pass"), report.Pass), color.FgGreen)
	PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", utils.I118Prt.Sprintf("fail"), report.Fail), color.FgRed)
	PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", utils.I118Prt.Sprintf("skip"), report.Skip), color.FgYellow)

	for _, cs := range report.Cases {
		str := "\n %s %s \n"
		status := cs.Status.String()
		statusColor := coloredStatus(status)

		logs = append(logs, fmt.Sprintf(str, status, cs.Path))
		Printt(fmt.Sprintf(str, statusColor, cs.Path))

		if len(cs.Steps) > 0 {
			count := 0
			for _, step := range cs.Steps {
				if count > 0 { // 空行
					PrintAndLog(&logs, "")
				}

				str := "  %s%d: %s   %s"
				status := utils.BoolToPass(step.Status)
				statusColor := coloredStatus(status)

				logs = append(logs, fmt.Sprintf(str, utils.I118Prt.Sprintf("step"), step.Numb, status, step.Name))
				Printt(fmt.Sprintf(str, utils.I118Prt.Sprintf("step"), step.Numb, statusColor, step.Name+"\n"))

				count1 := 0
				for _, cp := range step.CheckPoints {
					if count1 > 0 { // 空行
						PrintAndLog(&logs, "")
					}

					cpStatus := utils.BoolToPass(step.Status)
					cpStatusColored := coloredStatus(cpStatus)
					logs = append(logs, fmt.Sprintf("    %s%d: %s\n", utils.I118Prt.Sprintf("checkpoint"), cp.Numb,
						utils.BoolToPass(cp.Status)))
					Printt(fmt.Sprintf("    %s%d: %s\n", utils.I118Prt.Sprintf("checkpoint"), cp.Numb, cpStatusColored))

					PrintAndLog(&logs, fmt.Sprintf("      %s %s", utils.I118Prt.Sprintf("expect_result"), cp.Expect))
					PrintAndLog(&logs, fmt.Sprintf("      %s %s", utils.I118Prt.Sprintf("actual_result"), cp.Actual))

					count1++
				}

				count++
			}
		} else {
			PrintAndLog(&logs, "   "+utils.I118Prt.Sprintf("no_checkpoints"))
		}
	}

	utils.WriteFile(workDir+"/logs/result-"+utils.DateTimeStrLong(time.Now())+".txt", strings.Join(logs, "\n"))
}

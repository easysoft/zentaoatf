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
	PrintAndLogColorLn(&logs, fmt.Sprintf("%s: %d", utils.I118Prt.Sprintf("pass"), report.Pass), color.FgGreen)
	PrintAndLogColorLn(&logs, fmt.Sprintf("%s: %d", utils.I118Prt.Sprintf("fail"), report.Fail), color.FgRed)
	PrintAndLogColorLn(&logs, fmt.Sprintf("%s: %d", utils.I118Prt.Sprintf("skip"), report.Skip), color.FgYellow)

	for _, cs := range report.Cases {
		str := "\n%s %s"
		status := cs.Status.String()
		statusColor := colorStatus(status)

		logs = append(logs, fmt.Sprintf(str, status, cs.Path))
		fmt.Printf(str, statusColor, cs.Path)

		if len(cs.Steps) > 0 {
			count := 0
			for _, step := range cs.Steps {
				if count > 0 { // 空行
					PrintAndLog(&logs, "")
				}

				str := "  Step %d %s: %s"
				status := utils.BoolToPass(step.Status)
				statusColor := colorStatus(status)

				logs = append(logs, fmt.Sprintf(str, step.Numb, step.Name, status))
				fmt.Printf(str, step.Numb, step.Name, statusColor+"\n")

				count1 := 0
				for _, cp := range step.CheckPoints {
					if count1 > 0 { // 空行
						PrintAndLog(&logs, "")
					}

					PrintAndLog(&logs, fmt.Sprintf("    Checkpoint %d: %s", cp.Numb,
						utils.BoolToPass(cp.Status)))
					PrintAndLog(&logs, fmt.Sprintf("      Expect %s", cp.Expect))
					PrintAndLog(&logs, fmt.Sprintf("      Actual %s", cp.Actual))

					count1++
				}

				count++
			}
		} else {
			PrintAndLog(&logs, "   No check points")
		}
	}

	utils.WriteFile(workDir+"/logs/result-"+utils.DateTimeStrLong(time.Now())+".txt", strings.Join(logs, "\n"))
}

func colorStatus(status string) string {
	temp := strings.ToLower(status)

	switch temp {
	case "pass":
		return color.GreenString(status)
	case "fail":
		return color.RedString(status)
	case "skip":
		return color.YellowString(status)
	}

	return status
}

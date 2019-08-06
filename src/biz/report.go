package biz

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"os"
	"strings"
	"time"
)

func Print(report model.TestReport, workDir string) {
	startSec := time.Unix(report.StartTime, 0)
	endSec := time.Unix(report.EndTime, 0)

	logs := make([]string, 0)

	utils.PrintAndLog(&logs, utils.I118Prt.Sprintf("run_scripts", report.Path, report.Env))

	utils.PrintAndLog(&logs, utils.I118Prt.Sprintf("time_from_to",
		startSec.Format("2006-01-02 15:04:05"), endSec.Format("2006-01-02 15:04:05"), report.Duration))

	utils.PrintAndLog(&logs, fmt.Sprintf("%s: %d", utils.I118Prt.Sprintf("total"), report.Total))
	utils.PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", utils.I118Prt.Sprintf("pass"), report.Pass), color.FgGreen)
	utils.PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", utils.I118Prt.Sprintf("fail"), report.Fail), color.FgRed)
	utils.PrintAndLogColorLn(&logs, fmt.Sprintf("  %s: %d", utils.I118Prt.Sprintf("skip"), report.Skip), color.FgYellow)

	for _, cs := range report.Cases {
		str := "\n %s %s"
		status := cs.Status.String()
		statusColor := utils.ColoredStatus(status)

		logs = append(logs, fmt.Sprintf(str, status, cs.Path))
		utils.Printt(fmt.Sprintf(str+"\n", statusColor, cs.Path))

		if len(cs.Steps) > 0 {
			count := 0
			for _, step := range cs.Steps {
				if count > 0 { // 空行
					utils.PrintAndLog(&logs, "")
				}

				str := "  %s%d: %s   %s"
				status := utils.BoolToPass(step.Status)
				statusColor := utils.ColoredStatus(status)

				logs = append(logs, fmt.Sprintf(str, utils.I118Prt.Sprintf("step"), step.Numb, status, step.Name))
				utils.Printt(fmt.Sprintf(str, utils.I118Prt.Sprintf("step"), step.Numb, statusColor, step.Name+"\n"))

				count1 := 0
				for _, cp := range step.CheckPoints {
					if count1 > 0 { // 空行
						utils.PrintAndLog(&logs, "")
					}

					cpStatus := utils.BoolToPass(step.Status)
					cpStatusColored := utils.ColoredStatus(cpStatus)
					logs = append(logs, fmt.Sprintf("    %s%d: %s", utils.I118Prt.Sprintf("checkpoint"), cp.Numb,
						utils.BoolToPass(cp.Status)))
					utils.Printt(fmt.Sprintf("    %s%d: %s\n", utils.I118Prt.Sprintf("checkpoint"), cp.Numb, cpStatusColored))

					utils.PrintAndLog(&logs, fmt.Sprintf("      %s %s", utils.I118Prt.Sprintf("expect_result"), cp.Expect))
					utils.PrintAndLog(&logs, fmt.Sprintf("      %s %s", utils.I118Prt.Sprintf("actual_result"), cp.Actual))

					count1++
				}

				count++
			}
		} else {
			utils.PrintAndLog(&logs, "   "+utils.I118Prt.Sprintf("no_checkpoints"))
		}
	}

	utils.WriteFile(workDir+string(os.PathSeparator)+utils.LogDir+utils.RunDir+"result.txt", strings.Join(logs, "\n"))
}

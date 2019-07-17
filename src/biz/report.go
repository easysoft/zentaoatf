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

	PrintAndLog(&logs, fmt.Sprintf("Run scripts in folder \"%s\" on %s OS\n",
		report.Path, report.Env))

	PrintAndLog(&logs, fmt.Sprintf("From %s to %s, duration %d sec",
		startSec.Format("2006-01-02 15:04:05"), endSec.Format("2006-01-02 15:04:05"), report.Duration))

	PrintAndLog(&logs, fmt.Sprintf("Total: %d \n %s\n %s\n %s",
		report.Total,
		color.GreenString("Pass: %d", report.Pass),
		color.RedString("Fail: %d", report.Fail),
		color.BlueString("Skip: %d", report.Skip)))

	for _, cs := range report.Cases {
		PrintAndLog(&logs, fmt.Sprintf("\n%s %s", colorStatus(cs.Status.String()), cs.Path))

		if len(cs.Steps) > 0 {
			count := 0
			for _, step := range cs.Steps {
				if count > 0 { // 空行
					PrintAndLog(&logs, "")
				}

				PrintAndLog(&logs, fmt.Sprintf("  Step %d %s: %s", step.Numb, step.Name,
					colorStatus(utils.BoolToPass(step.Status))))

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

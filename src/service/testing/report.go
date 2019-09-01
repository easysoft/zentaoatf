package testingService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"time"
)

func Report(report model.TestReport) {
	startSec := time.Unix(report.StartTime, 0)
	endSec := time.Unix(report.EndTime, 0)

	logUtils.TraceAndResult(i118Utils.I118Prt.Sprintf("run_scripts", vari.WorkDir, report.Env))
	logUtils.TraceAndResult(i118Utils.I118Prt.Sprintf("time_from_to",
		startSec.Format("2006-01-02 15:04:05"), endSec.Format("2006-01-02 15:04:05"), report.Duration))

	logUtils.TraceAndResult(fmt.Sprintf("%s: %d", i118Utils.I118Prt.Sprintf("total"), report.Total))
	logUtils.TraceAndResult(fmt.Sprintf("  %s: %d, %s: %d, %s: %d",
		i118Utils.I118Prt.Sprintf("pass"), report.Pass,
		i118Utils.I118Prt.Sprintf("fail"), report.Fail,
		i118Utils.I118Prt.Sprintf("skip"), report.Skip))

	for idx, cs := range report.Cases {
		status := i118Utils.I118Prt.Sprintf(cs.Status)

		logUtils.TraceAndResult(fmt.Sprintf("\n%d. %s %s, %s (%d/%d)",
			cs.Id, cs.Title, status, cs.Path, idx+1, len(report.Cases)))

		if len(cs.Steps) > 0 {
			count := 0
			for _, step := range cs.Steps {
				if count > 0 { // 空行
					logUtils.PrintAndLog("")
				}

				status := i118Utils.I118Prt.Sprintf(commonUtils.BoolToPass(step.Status))
				logUtils.TraceAndResult(fmt.Sprintf("  [Step%d]: %s %s", step.Id, step.Name, status))

				count1 := 0
				for _, cp := range step.CheckPoints {
					if count1 > 0 { // 空行
						logUtils.PrintAndLog("")
					}

					//cpStatus := commonUtils.BoolToPass(step.Status)
					logUtils.TraceAndResult(fmt.Sprintf("    [Expect] %s", cp.Expect))
					logUtils.TraceAndResult(fmt.Sprintf("    [Actual] %s", cp.Actual))

					count1++
				}

				count++
			}
		} else {
			logUtils.TraceAndResult("   " + i118Utils.I118Prt.Sprintf("no_checkpoints"))
		}
	}

	json, _ := json.Marshal(report)
	fileUtils.WriteFile(constant.LogDir+vari.RunDir+"result.json", string(json))
}

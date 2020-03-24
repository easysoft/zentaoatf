package testingService

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"strconv"
	"strings"
	"time"
)

func GenUnitTestReport(cases []model.UnitTestCase, classNameMaxWidth int) model.UnitTestReport {
	logUtils.InitLogger()
	report := model.UnitTestReport{Env: commonUtils.GetOs(), Pass: 0, Fail: 0, Total: 0}

	failedCount := 0
	failedCaseLines := make([]string, 0)
	failedCaseLinesDesc := make([]string, 0)

	for idx, cs := range cases {
		cases[idx].Classname = stringUtils.AddPostfix(cs.Classname, classNameMaxWidth, " ")

		if cs.Failure != nil {
			report.Fail++

			if failedCount > 0 { // 换行
				failedCaseLinesDesc = append(failedCaseLinesDesc, "")
			}
			className := cases[idx].Classname

			line := fmt.Sprintf("[%s] %d.%s", className, cs.Id, cs.Name)
			failedCaseLines = append(failedCaseLines, line)

			failedCaseLinesDesc = append(failedCaseLinesDesc, line)
			failDesc := fmt.Sprintf("   %s - %s", cs.Failure.Type, cs.Failure.Desc)
			failedCaseLinesDesc = append(failedCaseLinesDesc, failDesc)
		} else {
			report.Pass++
		}
		report.Total++
	}
	report.Cases = cases

	postFix := ":"
	if len(cases) == 0 {
		postFix = "."
	}

	logUtils.ScreenAndResult(time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("found_scripts", color.CyanString(strconv.Itoa(len(cases)))) + postFix)

	if report.Total == 0 {
		return report
	}

	width := strconv.Itoa(len(strconv.Itoa(report.Total)))
	for idx, cs := range cases {
		statusColor := logUtils.ColoredStatus(cs.Status)

		format := "(%" + width + "d/%d) %s [%s] [%" + width + "d. %s] (%.3fs)"
		logUtils.Screen(fmt.Sprintf(format, idx+1, report.Total, statusColor, cs.Classname, cs.Id, cs.Name, cs.Time))
		logUtils.Result(fmt.Sprintf(format, idx+1, report.Total,
			i118Utils.I118Prt.Sprintf(cs.Status), cs.Classname, cs.Id, cs.Name, cs.Time))
	}

	if report.Fail > 0 {
		logUtils.ScreenAndResult("\n" + i118Utils.I118Prt.Sprintf("failed_scripts"))
		logUtils.Screen(strings.Join(failedCaseLines, "\n"))
		logUtils.Result(strings.Join(failedCaseLinesDesc, "\n"))
	}

	secTag := ""
	if vari.Config.Language == "en" && report.Duration > 1 {
		secTag = "s"
	}

	fmtStr := "%d(%.1f%%) %s"
	passStr := fmt.Sprintf(fmtStr, report.Pass, float32(report.Pass*100/report.Total), i118Utils.I118Prt.Sprintf("pass"))
	failStr := fmt.Sprintf(fmtStr, report.Fail, float32(report.Fail*100/report.Total), i118Utils.I118Prt.Sprintf("fail"))
	skipStr := fmt.Sprintf(fmtStr, report.Skip, float32(report.Skip*100/report.Total), i118Utils.I118Prt.Sprintf("skip"))

	// 输出到文件
	logUtils.Result("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("run_scripts",
			report.Total, report.Duration, secTag,
			passStr, failStr, skipStr,
			vari.LogDir+"result.txt",
		))

	// 输出到屏幕
	logUtils.Screen("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("run_scripts",
			report.Total, report.Duration, secTag,
			color.GreenString(passStr), color.RedString(failStr), color.YellowString(skipStr),
			vari.LogDir+"result.txt",
		))

	json, _ := json.Marshal(report)
	fileUtils.WriteFile(vari.LogDir+"result.json", string(json))

	return report
}

func SubmitUnitTestReport(report model.UnitTestReport) {

}

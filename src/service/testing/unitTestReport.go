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

func GenUnitTestReport(cases []model.UnitResult, classNameMaxWidth int, duration float32) model.TestReport {
	logUtils.InitLogger()
	report := model.TestReport{Env: commonUtils.GetOs(), Pass: 0, Fail: 0, Total: 0}
	report.TestType = "unit"
	report.TestFrame = vari.UnitTestType

	failedCount := 0
	failedCaseLines := make([]string, 0)
	failedCaseLinesDesc := make([]string, 0)

	for idx, cs := range cases {
		if cs.Failure != nil {
			report.Fail++

			if failedCount > 0 { // 换行
				failedCaseLinesDesc = append(failedCaseLinesDesc, "")
			}
			className := cases[idx].TestSuite

			line := fmt.Sprintf("[%s] %d.%s", className, cs.Id, cs.Title)
			failedCaseLines = append(failedCaseLines, line)

			failedCaseLinesDesc = append(failedCaseLinesDesc, line)
			failDesc := fmt.Sprintf("   %s - %s", cs.Failure.Type, cs.Failure.Desc)
			failedCaseLinesDesc = append(failedCaseLinesDesc, failDesc)
		} else {
			report.Pass++
		}
		report.Total++

		if duration == 0 {
			if report.StartTime == 0 || cs.StartTime < report.StartTime {
				report.StartTime = cs.StartTime
			}
			if cs.EndTime > report.EndTime {
				report.EndTime = cs.EndTime
			}
		}
	}
	report.UnitResult = cases
	if duration == 0 {
		report.Duration = report.EndTime - report.StartTime
	} else {
		report.Duration = int64(duration)
	}

	postFix := ":"
	if len(cases) == 0 {
		postFix = "."
	}

	logUtils.Result("\n" + logUtils.GetWholeLine(time.Now().Format("2006-01-02 15:04:05")+" "+
		i118Utils.Sprintf("found_scripts", strconv.Itoa(len(cases)))+postFix, "="))
	logUtils.Screen("\n" + logUtils.GetWholeLine(time.Now().Format("2006-01-02 15:04:05")+" "+
		i118Utils.Sprintf("found_scripts", color.CyanString(strconv.Itoa(len(cases))))+postFix, "="))

	if report.Total == 0 {
		return report
	}

	width := strconv.Itoa(len(strconv.Itoa(report.Total)))
	for idx, cs := range cases {
		statusColor := logUtils.ColoredStatus(cs.Status)
		testSuite := stringUtils.AddPostfix(cs.TestSuite, classNameMaxWidth, " ")

		format := "(%" + width + "d/%d) %s [%s] [%" + width + "d. %s] (%.3fs)"
		logUtils.Screen(fmt.Sprintf(format, idx+1, report.Total, statusColor, testSuite, cs.Id, cs.Title, cs.Duration))
		logUtils.Result(fmt.Sprintf(format, idx+1, report.Total,
			i118Utils.Sprintf(cs.Status), testSuite, cs.Id, cs.Title, cs.Duration))
	}

	if report.Fail > 0 {
		logUtils.ScreenAndResult("\n" + i118Utils.Sprintf("failed_scripts"))
		logUtils.Screen(strings.Join(failedCaseLines, "\n"))
		logUtils.Result(strings.Join(failedCaseLinesDesc, "\n"))
	}

	secTag := ""
	if vari.Config.Language == "en" && report.Duration > 1 {
		secTag = "s"
	}

	fmtStr := "%d(%.1f%%) %s"
	passStr := fmt.Sprintf(fmtStr, report.Pass, float32(report.Pass*100/report.Total), i118Utils.Sprintf("pass"))
	failStr := fmt.Sprintf(fmtStr, report.Fail, float32(report.Fail*100/report.Total), i118Utils.Sprintf("fail"))
	skipStr := fmt.Sprintf(fmtStr, report.Skip, float32(report.Skip*100/report.Total), i118Utils.Sprintf("skip"))

	// 输出到文件
	logUtils.Result("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.Sprintf("run_scripts",
			report.Total, report.Duration, secTag,
			passStr, failStr, skipStr,
			" "+vari.LogDir+"result.txt ",
		))

	// 输出到屏幕
	logUtils.Screen("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.Sprintf("run_scripts",
			report.Total, report.Duration, secTag,
			color.GreenString(passStr), color.RedString(failStr), color.YellowString(skipStr),
			" "+vari.LogDir+"result.txt ",
		))

	json, _ := json.Marshal(report)
	fileUtils.WriteFile(vari.LogDir+"result.json", string(json))

	return report
}

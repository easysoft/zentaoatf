package execHelper

import (
	"fmt"
	"path/filepath"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	"github.com/fatih/color"
)

func GenStatusTxt(status commConsts.ResultStatus) (txt string) {
	txt = i118Utils.Sprintf(string(status))
	if commConsts.ExecFrom == commConsts.FromCmd {
		if status == commConsts.FAIL {
			txt = color.New(color.FgHiRed, color.Bold).Sprint(status)
		} else if status == commConsts.PASS {
			txt = color.New(color.FgHiGreen, color.Bold).Sprint(status)
		} else {
			txt = color.New(color.FgHiYellow, color.Bold).Sprint(status)
		}
	}

	return
}

func ZipDir(req serverDomain.TestSet) {
	if req.ZipDir == "" {
		return
	}

	zipFile := filepath.Join(commConsts.ExecLogDir, commConsts.ResultZip)
	fileUtils.ZipDir(zipFile, req.ZipDir)
}

func GenZTFReport(report commDomain.ZtfReport) (failedCount int, failedCaseLinesWithCheckpoint []string) {
	failedCaseLinesWithCheckpoint = make([]string, 0)

	for index, csResult := range report.FuncResult {
		if report.ProductId == 0 && csResult.ProductId > 0 {
			report.ProductId = csResult.ProductId
		}

		if csResult.Status == "fail" {
			if failedCount > 0 {
				failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
			}
			failedCount++

			path := csResult.Path
			relativePath := strings.TrimPrefix(path, commConsts.WorkDir)

			prefix := i118Utils.Sprintf("test_case_prefix", index+1)
			line := fmt.Sprintf("%s [%s] [%s]", prefix, relativePath, csResult.Title)
			failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, line)

			appendFailedStepResult(csResult, &failedCaseLinesWithCheckpoint)
		}
	}

	return
}

func GenUnitReport(cases []commDomain.UnitResult, report *commDomain.ZtfReport, duration float32) (failedCaseLines, failedCaseLinesDesc []string) {
	failedCaseLines = make([]string, 0)
	failedCaseLinesDesc = make([]string, 0)

	for idx, cs := range cases {
		if cs.Failure != nil {
			report.Fail++

			className := cases[idx].TestSuite

			line := fmt.Sprintf("[%s] %s", className, cs.Title)
			failedCaseLines = append(failedCaseLines, line)

			failedCaseLinesDesc = append(failedCaseLinesDesc, line)
			failDesc := fmt.Sprintf("   %s - %s", cs.Failure.Type, cs.Failure.Desc)
			failedCaseLinesDesc = append(failedCaseLinesDesc, failDesc)
		} else {
			report.Pass++
		}

		report.Total++

		if cs.EndTime > report.EndTime {
			report.EndTime = cs.EndTime
		}
	}

	report.UnitResult = cases
	if duration == 0 {
		report.Duration = report.EndTime - report.StartTime
	} else {
		report.Duration = int64(duration)
	}

	return
}

func GenRunResult(report commDomain.ZtfReport) (result, resultClient string) {
	fmtStr := "%s%s%d(%.1f%%)"
	passRate := 0
	failRate := 0
	skipRate := 0
	if report.Total > 0 {
		passRate = report.Pass * 100 / report.Total
		failRate = report.Fail * 100 / report.Total
		skipRate = report.Skip * 100 / report.Total
	}

	passStr := fmt.Sprintf(fmtStr, i118Utils.Sprintf("pass_num"), i118Utils.Sprintf("colon"), report.Pass, float32(passRate))
	failStr := fmt.Sprintf(fmtStr, i118Utils.Sprintf("fail_num"), i118Utils.Sprintf("colon"), report.Fail, float32(failRate))
	skipStr := fmt.Sprintf(fmtStr, i118Utils.Sprintf("skip_num"), i118Utils.Sprintf("colon"), report.Skip, float32(skipRate))

	if commConsts.ExecFrom == commConsts.FromCmd {
		passStr = fmt.Sprintf(fmtStr, color.New(color.FgHiGreen, color.Bold).Sprint(i118Utils.Sprintf("pass_num")), i118Utils.Sprintf("colon"), report.Pass, float32(passRate))
		failStr = fmt.Sprintf(fmtStr, color.New(color.FgHiRed, color.Bold).Sprint(i118Utils.Sprintf("fail_num")), i118Utils.Sprintf("colon"), report.Fail, float32(failRate))
		skipStr = fmt.Sprintf(fmtStr, color.New(color.FgHiYellow, color.Bold).Sprint(i118Utils.Sprintf("skip_num")), i118Utils.Sprintf("colon"), report.Skip, float32(skipRate))
	}

	if commConsts.ExecFrom == commConsts.FromClient {
		resultClient = i118Utils.Sprintf("run_result", report.Total, report.Duration,
			fmt.Sprintf(`<span class="result-pass">%s</span>`, passStr),
			fmt.Sprintf(`<span class="result-fail">%s</span>`, failStr),
			fmt.Sprintf(`<span class="result-skip">%s</span>`, skipStr),
		)
	}

	result = i118Utils.Sprintf("run_result",
		report.Total, report.Duration,
		passStr, failStr, skipStr,
	)

	return
}

func SetRunningIfNeed(ch chan int) {
	select {
	case _, ok := <-ch:
		if !ok {
			SetRunning(false)
			return
		}
	default:
	}
}

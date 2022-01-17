package reportUtils

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GenZTFTestReport(report commDomain.ZtfReport, pathMaxWidth int, projectPath string) {
	if len(report.FuncResult) == 0 {
		return
	}

	// print failed case
	failedCount := 0
	failedCaseLines := make([]string, 0)
	failedCaseLinesWithCheckpoint := make([]string, 0)

	for _, cs := range report.FuncResult {
		if cs.Status == "fail" {
			if failedCount > 0 {
				failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
			}
			failedCount++

			path := cs.Path
			lent := runewidth.StringWidth(path)

			if pathMaxWidth > lent {
				postFix := strings.Repeat(" ", pathMaxWidth-lent)
				path += postFix
			}

			line := fmt.Sprintf("[%s] %d.%s", cs.Path, cs.Id, cs.Title)
			failedCaseLines = append(failedCaseLines, line)
			failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, line)

			if len(cs.Steps) > 0 {
				stepNumb := 0
				for _, step := range cs.Steps {
					if step.Status {
						continue
					}

					if stepNumb > 0 {
						failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
					}
					stepNumb++

					step.Id = strings.TrimRight(step.Id, ".")
					status := i118Utils.Sprintf(stringUtils.BoolToPass(step.Status))
					failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("Step %s: %s", step.Id, status))

					for idx1, cp := range step.CheckPoints {
						//cpStatus := commonUtils.BoolToPass(step.Status)
						failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("[Expect] %s", cp.Expect))
						failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, fmt.Sprintf("[Actual] %s", cp.Actual))

						if idx1 < len(step.CheckPoints)-1 {
							failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "")
						}
					}
				}
			} else {
				failedCaseLinesWithCheckpoint = append(failedCaseLinesWithCheckpoint, "   "+i118Utils.Sprintf("no_checkpoints"))
			}
		}
	}
	if failedCount > 0 {
		logUtils.ScreenAndResult("\n" + i118Utils.Sprintf("failed_scripts"))
		logUtils.Screen(strings.Join(failedCaseLines, "\n"))
		logUtils.Result(strings.Join(failedCaseLinesWithCheckpoint, "\n"))
	}

	secTag := ""
	if commConsts.Language == "en" && report.Duration > 1 {
		secTag = "s"
	}

	// 生成统计行
	fmtStr := "%d(%.1f%%) %s"
	passStr := fmt.Sprintf(fmtStr, report.Pass, float32(report.Pass*100/report.Total), i118Utils.Sprintf("pass"))
	failStr := fmt.Sprintf(fmtStr, report.Fail, float32(report.Fail*100/report.Total), i118Utils.Sprintf("fail"))
	skipStr := fmt.Sprintf(fmtStr, report.Skip, float32(report.Skip*100/report.Total), i118Utils.Sprintf("skip"))

	logFile := filepath.Join(projectPath, commConsts.LogDir, "result.txt")

	// 打印到结果文件
	logUtils.Result("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.Sprintf("run_scripts",
			report.Total, report.Duration, secTag,
			passStr, failStr, skipStr,
			" "+logFile,
		))

	//report.ProductId, _ = strconv.Atoi(vari.ProductId)
	json, _ := json.Marshal(report)
	fileUtils.WriteFile(logFile, string(json))
}

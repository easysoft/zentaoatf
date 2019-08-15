package testingService

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/src/model"
	scriptService "github.com/easysoft/zentaoatf/src/service/script"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
	"strings"
)

func GetTestTestReportForSubmit(assert string, date string) model.TestReport {
	mode, name := scriptService.GetRunModeAndName(assert)
	resultPath := vari.Prefer.WorkDir + constant.LogDir + scriptService.LogFolder(mode, name, date) +
		string(os.PathSeparator) + "result.json"

	content := fileUtils.ReadFile(resultPath)
	content = strings.Replace(content, "\n", "", -1)

	var report model.TestReport
	json.Unmarshal([]byte(content), &report)

	return report
}

func SaveTestTestReportAfterSubmit(assert string, date string, content string) {
	mode, name := scriptService.GetRunModeAndName(assert)
	resultPath := vari.Prefer.WorkDir + constant.LogDir + scriptService.LogFolder(mode, name, date) +
		string(os.PathSeparator) + "result.json"

	fileUtils.WriteFile(resultPath, content)
}

package action

import (
	"github.com/easysoft/zentaoatf/src/ui/page"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
)

func CommitBug(files []string, caseId string) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		configUtils.ConfigForDir(&resultDir, "result")
	}
	resultDir = fileUtils.UpdateDir(resultDir)

	if caseId == "" {
		configUtils.ConfigForInt(&caseId, "test_case")
	}

	page.CuiReportBug("logs/2019-08-30T130258/", "1")
}

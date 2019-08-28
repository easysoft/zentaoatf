package action

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
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

	if caseId == "" {
		configUtils.ConfigForInt(&caseId, "test_case")
	}

	resultDir = fileUtils.UpdateDir(resultDir)

	bug, stepIds := zentaoService.GenBug(resultDir, caseId)
	zentaoService.CommitBug(bug, stepIds)
}

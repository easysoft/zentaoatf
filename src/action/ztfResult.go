package action

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
)

func CommitZTFTestResult(files []string, noNeedConfirm bool) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		stdinUtils.InputForDir(&resultDir, "", "result")
	}

	resultDir = fileUtils.UpdateDir(resultDir)
	zentaoService.CommitZTFTestResult(resultDir, noNeedConfirm)

}

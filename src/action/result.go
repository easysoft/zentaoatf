package action

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
)

func CommitResult(files []string) {
	var resultDir string
	if len(files) > 0 {
		resultDir = files[0]
	} else {
		configUtils.ConfigForDir(&resultDir, "result")
	}

	resultDir = fileUtils.UpdateDir(resultDir)
	zentaoService.CommitResult(resultDir)

}

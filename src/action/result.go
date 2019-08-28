package action

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
)

func CommitResult(resultDir string) {
	resultDir = fileUtils.UpdateDir(resultDir)
	zentaoService.CommitResult(resultDir)
}

package action

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/ui/page"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
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

	vari.CurrBug = zentaoService.PrepareBug(resultDir, caseId)
	page.Cui("bug")
}

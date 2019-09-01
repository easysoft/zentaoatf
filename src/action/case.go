package action

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
)

func CommitCases(files []string) {
	for _, file := range files {
		id, _, title := zentaoUtils.GetCaseInfo(file)
		zentaoService.CommitCase(id, title)
	}
}

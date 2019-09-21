package action

import (
	scriptService "github.com/easysoft/zentaoatf/src/service/script"
	assertUtils "github.com/easysoft/zentaoatf/src/utils/assert"
)

func View(files []string, keywords string) {
	cases := assertUtils.GetCaseByDirAndFile(files)

	scriptService.View(cases, keywords)
}

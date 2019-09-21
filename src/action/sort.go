package action

import (
	scriptService "github.com/easysoft/zentaoatf/src/service/script"
	assertUtils "github.com/easysoft/zentaoatf/src/utils/assert"
)

func Sort(files []string) {
	cases := assertUtils.GetCaseByDirAndFile(files)

	scriptService.Sort(cases)
}

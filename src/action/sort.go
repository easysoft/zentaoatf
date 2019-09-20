package action

import scriptService "github.com/easysoft/zentaoatf/src/service/script"

func Sort(files []string) {
	cases := scriptService.GetCaseByDirAndFile(files)

	scriptService.Sort(cases)
}

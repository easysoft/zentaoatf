package action

import scriptService "github.com/easysoft/zentaoatf/src/service/script"

func List(files []string, keywords string) {
	cases := scriptService.GetCaseByDirAndFile(files)

	scriptService.List(cases, keywords)
}

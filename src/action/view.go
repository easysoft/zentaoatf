package action

import scriptService "github.com/easysoft/zentaoatf/src/service/script"

func View(files []string, keywords string) {
	cases := make([]string, 0)

	for _, file := range files {
		scriptService.GetAllScriptsInDir(file, &cases)
	}

	scriptService.View(cases, keywords)
}

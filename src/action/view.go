package action

import "github.com/easysoft/zentaoatf/src/biz"

func View(scriptDir string, files []string, langType string) {
	biz.View(scriptDir, files, langType)
}

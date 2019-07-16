package action

import "github.com/easysoft/zentaoatf/src/biz"

func View(scriptDir string, langType string, files []string) {
	biz.View(scriptDir, langType, files)
}

package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
)

func View(scriptDir string, files []string, langType string) {
	scriptService.View(scriptDir, files, langType)
}

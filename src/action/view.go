package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
)

func View(dir string, files []string, langType string) {
	scriptService.View(dir, files, langType)
}

package action

import (
	"github.com/easysoft/zentaoatf/src/service/script"
)

func List(scriptDir string, langType string) {
	scriptService.List(scriptDir, langType)
}

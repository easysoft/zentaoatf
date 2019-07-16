package biz

import (
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
)

func View(scriptDir string, fileNames []string, langType string) {
	var files []string
	if fileNames != nil && len(fileNames) > 0 {
		files, _ = utils.GetSpecifiedFiles(scriptDir, fileNames)
	} else {
		files, _ = utils.GetAllFiles(scriptDir, langType)
	}

	for _, file := range files {
		script.Detail(file)
	}

}

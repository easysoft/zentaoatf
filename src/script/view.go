package script

import (
	"github.com/easysoft/zentaoatf/src/utils/file"
)

func View(scriptDir string, fileNames []string, langType string) {
	files := make([]string, 0)
	if fileNames != nil && len(fileNames) > 0 {
		files, _ = fileUtils.GetSpecifiedFiles(scriptDir, fileNames)
	} else {
		fileUtils.GetAllFiles(scriptDir, langType, &files)
	}

	for _, file := range files {
		Detail(file)
	}

}

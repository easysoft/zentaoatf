package biz

import (
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
	"os"
	"strings"
)

func View(scriptDir string, langType string, fileNames []string) {

	files, _ := utils.GetAllFiles(scriptDir, langType)

	if files != nil && len(fileNames) > 0 {
		sep := string(os.PathSeparator)
		for _, name := range fileNames {
			file := name
			if strings.Index(file, sep) == -1 {
				file = scriptDir + sep + file
			}
			script.Detail(file)
		}
	}

}

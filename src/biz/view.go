package biz

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"os"
	"strings"
)

func View(scriptDir string, langType string, fileNames []string) {

	files, _ := utils.GetAllFiles(scriptDir, langType)

	fmt.Println(fileNames)

	if files != nil && len(fileNames) > 0 {
		sep := string(os.PathSeparator)
		for _, name := range fileNames {
			if strings.Index(name, sep) == -1 {
				name = scriptDir + sep + name
			}
			ReadFile(name)
		}
	}

}

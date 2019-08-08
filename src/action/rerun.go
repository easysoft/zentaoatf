package action

import (
	. "github.com/easysoft/zentaoatf/src/utils"
	"os"
)

func Rerun(resultFile string) {
	files, scriptDir, _, _ := GetFailedFiles(resultFile)

	if !PathEndWithSeparator(scriptDir) {
		scriptDir += string(os.PathSeparator)
	}

	Run(scriptDir, files, "")
}

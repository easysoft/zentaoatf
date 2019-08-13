package action

import (
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"os"
)

func Rerun(resultFile string) {
	files, scriptDir, _, _ := fileUtils.GetFailedFiles(resultFile)

	if !commonUtils.PathEndWithSeparator(scriptDir) {
		scriptDir += string(os.PathSeparator)
	}

	Run(scriptDir, files, "")
}

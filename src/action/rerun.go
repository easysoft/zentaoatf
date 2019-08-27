package action

import (
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"os"
)

func Rerun(resultFile string) {
	files, dir := fileUtils.GetFailedFilesFromTestResult(resultFile)

	if !commonUtils.PathEndWithSeparator(dir) {
		dir += string(os.PathSeparator)
	}

	Run(dir, files)
}

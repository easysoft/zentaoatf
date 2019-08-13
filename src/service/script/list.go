package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/file"
)

func List(scriptDir string, langType string) {
	files := make([]string, 0)
	fileUtils.GetAllFiles(scriptDir, langType, &files)

	fmt.Printf("Totally %d test cases \n", len(files))

	for _, file := range files {
		Summary(file)
	}
}

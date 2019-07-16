package biz

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
)

func List(scriptDir string, langType string) {
	files, _ := utils.GetAllFiles(scriptDir, langType)

	fmt.Printf("Totally %d test cases \n", len(files))

	for _, file := range files {
		ReadFile(file)
	}
}

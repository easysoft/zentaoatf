package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/fatih/color"
	"regexp"
)

func List(scriptDir string, langType string) {
	files := make([]string, 0)
	fileUtils.GetAllFiles(scriptDir, langType, &files)

	fmt.Printf("Totally %d test cases \n", len(files))

	for _, file := range files {
		Summary(file)
	}
}

func Summary(file string) {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(`<<TC[\S\s]*caseId:([^\n]*)\n+title:([^\n]*)\n`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		title := commonUtils.RemoveBlankLine(arr[2])

		fmt.Printf("%s %s \n", color.CyanString(caseId), title)
	}
}

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

func Detail(file string) {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(
		`<<TC[\S\s]*caseId:([^\n]*)\n+title:([^\n]*)\n+steps:([\S\s]*)\n+expects:([\S\s]*?)\n+(readme:|TC;)`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		title := commonUtils.RemoveBlankLine(arr[2])
		steps := commonUtils.RemoveBlankLine(arr[3])
		expects := commonUtils.RemoveBlankLine(arr[4])

		fmt.Printf("%s %s \n", color.CyanString(caseId), title)
		fmt.Printf("%s \n", steps)
		fmt.Printf("%s \n\n", expects)
	}
}

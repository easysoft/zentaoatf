package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/fatih/color"
	"regexp"
)

func List(dir string, langType string) {
	files := make([]string, 0)
	fileUtils.GetAllFilesInDir(dir, &files)

	fmt.Printf("Totally %d test cases \n", len(files))

	for idx, file := range files {
		Summary(file, idx)
	}
}

func Summary(file string, inx int) {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(`<<<TC[\S\s]*caseId:([^\n]*)(?:[\S\s]+?)\n+title:([^\n]*)\n`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		title := commonUtils.RemoveBlankLine(arr[2])

		fmt.Printf("%d %s %s \n", inx+1, color.CyanString("tc-%s", caseId), title)
	}
}

func View(dir string, fileNames []string, langType string) {
	files := make([]string, 0)
	if fileNames != nil && len(fileNames) > 0 {
		files, _ = fileUtils.GetSpecifiedFilesInWorkDir(fileNames)
	} else {
		fileUtils.GetAllFilesInDir(dir, &files)
	}

	for _, file := range files {
		if fileNames != nil && len(fileNames) > 0 {
			Detail(file)
		} else {
			Brief(file)
		}

	}

}

func Brief(file string) {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(
		`<<<TC[\S\s]*` +
			`caseId:([^\n]*)\n+` +
			`caseIdInTask:([^\n]*)\n+` +
			`taskId:([^\n]*)\n+` +
			`title:([^\n]*)\n+` +
			`steps:([\S\s]*)\n` +
			`expects:([\S\s]*?)\n+` +
			`(readme:|TC)`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		//caseIdInTask := commonUtils.RemoveBlankLine(arr[2])
		//taskId := commonUtils.RemoveBlankLine(arr[3])

		title := commonUtils.RemoveBlankLine(arr[4])
		steps := commonUtils.RemoveBlankLine(arr[5])
		expects := commonUtils.RemoveBlankLine(arr[6])

		color.Cyan("\n%s %s \n", caseId, title)
		fmt.Printf("Steps: \n%s \n\n", steps)
		fmt.Printf("Expect Results: \n%s\n", expects)
	}
}

func Detail(file string) {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(
		`<<<TC[\S\s]*` +
			`caseId:([^\n]*)\n+` +
			`caseIdInTask:([^\n]*)\n+` +
			`taskId:([^\n]*)\n+` +
			`title:([^\n]*)\n+`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		title := commonUtils.RemoveBlankLine(arr[4])

		color.Cyan("\n%s %s \n", caseId, title)
		fmt.Printf("%s\n", content)
	}
}

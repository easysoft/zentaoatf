package script

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/fatih/color"
	"regexp"
)

func Summary(file string) {
	content := utils.ReadFile(file)

	myExp := regexp.MustCompile(`<<<TC[\S\s]*caseId:([^\n]*)\n+title:([^\n]*)\n`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := utils.RemoveBlankLine(arr[1])
		title := utils.RemoveBlankLine(arr[2])

		fmt.Printf("%s %s \n", color.CyanString(caseId), title)
	}
}

func Detail(file string) {
	content := utils.ReadFile(file)

	myExp := regexp.MustCompile(
		`<<<TC[\S\s]*caseId:([^\n]*)\n+title:([^\n]*)\n+steps:([\S\s]*)\n+expects:([\S\s]*?)\n+(readme:|TC;)`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := utils.RemoveBlankLine(arr[1])
		title := utils.RemoveBlankLine(arr[2])
		steps := utils.RemoveBlankLine(arr[3])
		expects := utils.RemoveBlankLine(arr[4])

		fmt.Printf("%s %s \n", color.CyanString(caseId), title)
		fmt.Printf("%s \n", steps)
		fmt.Printf("%s \n\n", expects)
	}
}

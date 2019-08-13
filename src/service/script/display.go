package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/common"
	file2 "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/fatih/color"
	"regexp"
)

func Summary(file string) {
	content := file2.ReadFile(file)

	myExp := regexp.MustCompile(`<<TC[\S\s]*caseId:([^\n]*)\n+title:([^\n]*)\n`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		title := commonUtils.RemoveBlankLine(arr[2])

		fmt.Printf("%s %s \n", color.CyanString(caseId), title)
	}
}

func Detail(file string) {
	content := file2.ReadFile(file)

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

package biz

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils"
	"regexp"
)

func ReadFile(file string) {
	content := utils.ReadFile(file)

	myExp := regexp.MustCompile(`<<<TC[\S\s]*caseId:([^\n]*)\n+title:([^\n]*)\n`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := utils.RemoveBlankLine(arr[1])
		title := utils.RemoveBlankLine(arr[2])

		fmt.Printf("%s %s \n", caseId, title)
	}
}

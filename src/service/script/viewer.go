package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/fatih/color"
	"regexp"
	"strings"
)

func List(cases []string, keywords string) {
	fmt.Printf("Totally %d test cases \n", len(cases))

	for idx, tc := range cases {
		Summary(tc, idx, keywords)
	}
}

func Summary(file string, inx int, keywords string) {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(`<<<TC[\S\s]*caseId:([^\n]*)(?:[\S\s]+?)\n+title:([^\n]*)\n`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		title := commonUtils.RemoveBlankLine(arr[2])

		if strings.Index(title, keywords) > -1 {
			fmt.Printf("%s %s \n", color.CyanString("%s", caseId), title)
		}
	}
}

func View(cases []string, keywords string) {
	for _, file := range cases {
		Brief(file, keywords)
		//Detail(file)
	}
}

func Brief(file string, keywords string) {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(
		`<<<TC[\S\s]*` +
			`caseId:([^\n]*)\n+` +
			`productId:([^\n]*)\n+` +
			`title:([^\n]*)\n+` +
			`steps:.*\n([\S\s]*)\n` +
			`expects:.*\n([\S\s]*?)\n` +
			`(readme:|TC)`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		_ = commonUtils.RemoveBlankLine(arr[2])

		title := commonUtils.RemoveBlankLine(arr[3])
		steps := arr[4]
		expects := commonUtils.RemoveBlankLine(arr[5])

		if strings.Index(title, keywords) > -1 {
			color.Cyan("\n%s %s \n", caseId, title)
			fmt.Printf("Steps: \n%s \n", steps)
			fmt.Printf("Expects: \n%s\n", expects)
		}
	}
}

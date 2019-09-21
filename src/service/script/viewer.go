package scriptUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/fatih/color"
	"regexp"
	"strconv"
	"strings"
)

func List(cases []string, keywords string) {
	keywords = strings.TrimSpace(keywords)

	count := 0
	for _, tc := range cases {
		if Summary(tc, keywords) {
			count++
		}
	}

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("total_test_case", count), -1)
}

func Summary(file string, keywords string) bool {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(`<<<TC[\S\s]*caseId:([^\n]*)(?:[\S\s]+?)\n+title:([^\n]*)\n`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		caseId := commonUtils.RemoveBlankLine(arr[1])
		title := commonUtils.RemoveBlankLine(arr[2])

		_, err := strconv.Atoi(keywords)
		var pass bool

		if err == nil && keywords == caseId { // int
			pass = true
		} else if strings.Index(title, keywords) > -1 {
			pass = true
		}

		if pass {
			fmt.Printf("%s %s \n", caseId, title)

			return true
		} else {
			return false
		}
	}
	return false
}

func View(cases []string, keywords string) {
	keywords = strings.TrimSpace(keywords)
	count := 0

	for _, file := range cases {
		if Brief(file, keywords) {
			count++
		}
		//Detail(file)
	}

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("total_test_case", count), -1)
}

func Brief(file string, keywords string) bool {
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

		_, err := strconv.Atoi(keywords)
		var pass bool

		if err == nil && keywords == caseId { // int
			pass = true
		} else if strings.Index(title, keywords) > -1 {
			pass = true
		}

		if pass {
			logUtils.PrintToStdOut(fmt.Sprintf("%s %s", caseId, title), color.FgCyan)
			fmt.Printf("Steps: \n%s \n", steps)
			fmt.Printf("Expects: \n%s\n\n", expects)

			return true
		} else {
			return false
		}
	}

	return false
}

package scriptUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
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

	logUtils.PrintToStdOut("\n"+i118Utils.I118Prt.Sprintf("total_test_case", count), -1)
}

func Summary(file string, keywords string) bool {
	pass, caseId, _, title := zentaoUtils.GetCaseInfo(file)

	if pass {
		_, err := strconv.Atoi(keywords)
		var pass bool

		if err == nil && keywords == strconv.Itoa(caseId) { // int
			pass = true
		} else if strings.Index(title, keywords) > -1 {
			pass = true
		}

		if pass {
			fmt.Printf("%d. %s \n", caseId, title)

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
			logUtils.PrintToStdOut("", -1)
			count++
		}
	}

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("total_test_case", count), -1)
}

func Brief(file string, keywords string) bool {
	content := fileUtils.ReadFile(file)

	myExp := regexp.MustCompile(
		`\[case\][\S\s]*` +
			`title=([^\n]*)\n+` +
			`cid=([^\n]*)\n+` +
			`pid=([^\n]*)\n+` +
			`([\S\s]*)\n*` +
			`\[esac\]`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		title := commonUtils.RemoveBlankLine(arr[1])
		caseId := commonUtils.RemoveBlankLine(arr[2])

		//productId := commonUtils.RemoveBlankLine(arr[3])
		steps := commonUtils.RemoveBlankLine(arr[4])

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

			return true
		} else {
			return false
		}
	}

	return false
}

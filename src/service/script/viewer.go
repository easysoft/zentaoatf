package scriptUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func List(cases []string, keywords string) {
	keywords = strings.TrimSpace(keywords)

	scriptArr := make([]model.CaseLog, 0)
	for _, tc := range cases {
		pass, cs := Summary(tc, keywords)
		if pass {
			scriptArr = append(scriptArr, cs)
		}
	}

	total := len(scriptArr)
	width := strconv.Itoa(len(strconv.Itoa(total)))

	logUtils.Screen(time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("found_scripts", total) + "\n")

	for idx, cs := range scriptArr {
		format := "(%" + width + "d/%d) [%s] %d.%s"
		logUtils.Screen(fmt.Sprintf(format, idx+1, total, cs.Path, cs.Id, cs.Title))
	}
}

func Summary(file string, keywords string) (bool, model.CaseLog) {
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
			//fmt.Printf("%d. %s \n", caseId, title)

			return true, model.CaseLog{Id: caseId, Title: title, Path: file}
		} else {
			return false, model.CaseLog{}
		}
	}
	return false, model.CaseLog{}
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

package action

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	scriptUtils "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/mattn/go-runewidth"
	"strconv"
	"strings"
	"time"
)

func List(files []string, keywords string) {
	var cases []string
	for _, v1 := range files {
		group := scriptUtils.LoadScriptByWorkspace(v1)
		for _, v2 := range group {
			cases = append(cases, v2)
		}
	}
	keywords = strings.TrimSpace(keywords)
	scriptArr := make([]commDomain.FuncResult, 0)
	pathMaxWidth := 0
	numbMaxWidth := 0
	for _, tc := range cases {
		pass, cs := SummaryObj(tc, keywords)
		if pass {
			scriptArr = append(scriptArr, cs)
		}

		if len(cs.Path) > pathMaxWidth {
			pathMaxWidth = len(cs.Path)
		}

		if len(tc) > numbMaxWidth {
			numbMaxWidth = len(strconv.Itoa(cs.Id))
		}
	}
	numbWidth := strconv.Itoa(numbMaxWidth)

	total := len(scriptArr)
	width := strconv.Itoa(len(strconv.Itoa(total)))

	logUtils.Info("\n" + time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.Sprintf("found_scripts", total, commConsts.WorkDir))

	for idx, cs := range scriptArr {
		path := cs.Path
		lent := runewidth.StringWidth(path)

		if pathMaxWidth > lent {
			postFix := strings.Repeat(" ", pathMaxWidth-lent)
			path += postFix
		}

		format := "(%" + width + "d/%d) [%s] [%" + numbWidth + "d. %s]"
		logUtils.Info(fmt.Sprintf(format, idx+1, total, path, cs.Id, cs.Title))
	}
	logUtils.Info("")
}

func SummaryObj(file string, keywords string) (bool, commDomain.FuncResult) {
	pass, caseId, _, title := scriptUtils.GetCaseInfo(file)

	if pass {
		_, err := strconv.Atoi(keywords)
		var pass bool

		if err == nil && keywords == strconv.Itoa(caseId) { // int
			pass = true
		} else if strings.Index(title, keywords) > -1 {
			pass = true
		}

		if pass {
			return true, commDomain.FuncResult{Id: caseId, Title: title, Path: file}
		} else {
			return false, commDomain.FuncResult{}
		}
	}
	return false, commDomain.FuncResult{}
}

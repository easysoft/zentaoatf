package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/pkg/consts"
	"strconv"
	"strings"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/mattn/go-runewidth"
)

func List(files []string, keywords string) {
	var cases []string
	for _, f := range files {
		group := scriptHelper.LoadScriptByWorkspace(f)
		for _, item := range group {
			cases = append(cases, item)
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

	logUtils.Info("\n" + time.Now().Format(consts.DateTimeFormat) + " " +
		i118Utils.Sprintf("found_scripts_no_ztf_dir", total, commConsts.WorkDir))

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
	pass, caseId, _, title, _ := scriptHelper.GetCaseInfo(file)

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

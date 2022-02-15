package action

import (
	"fmt"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	"github.com/aaronchen2k/deeptest/internal/command"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"strconv"
	"strings"
	"time"
)

func List(files []string, keywords string, actionModule *command.IndexModule) {
	cases := scriptUtils.LoadScriptByProject(files[0])
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

	logUtils.Info(time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.Sprintf("found_scripts", color.CyanString(strconv.Itoa(total))) + "\n")

	for idx, cs := range scriptArr {
		//format := "(%" + width + "d/%d) [%s] %d.%s"
		//logUtils.Screen(fmt.Sprintf(format, idx+1, total, cs.Path, cs.Id, cs.Title))

		path := cs.Path
		lent := runewidth.StringWidth(path)

		if pathMaxWidth > lent {
			postFix := strings.Repeat(" ", pathMaxWidth-lent)
			path += postFix
		}

		format := "(%" + width + "d/%d) [%s] [%" + numbWidth + "d. %s]"
		logUtils.Info(fmt.Sprintf(format, idx+1, total, path, cs.Id, cs.Title))
	}
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

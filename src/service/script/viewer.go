package scriptUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	langUtils "github.com/easysoft/zentaoatf/src/utils/lang"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func List(cases []string, keywords string) {
	keywords = strings.TrimSpace(keywords)

	scriptArr := make([]model.FuncResult, 0)

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

	logUtils.Screen(time.Now().Format("2006-01-02 15:04:05") + " " +
		i118Utils.I118Prt.Sprintf("found_scripts", color.CyanString(strconv.Itoa(total))) + "\n")

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
		logUtils.Screen(fmt.Sprintf(format, idx+1, total, path, cs.Id, cs.Title))
	}
}

func SummaryObj(file string, keywords string) (bool, model.FuncResult) {
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

			return true, model.FuncResult{Id: caseId, Title: title, Path: file}
		} else {
			return false, model.FuncResult{}
		}
	}
	return false, model.FuncResult{}
}

func View(cases []string, keywords string) {
	keywords = strings.TrimSpace(keywords)
	count := 0

	arrs := make([][]string, 0)
	for _, file := range cases {
		pass, arr := Brief(file, keywords)

		if pass {
			arrs = append(arrs, arr)
			count++
		}
	}

	total := len(arrs)

	logUtils.Screen(time.Now().Format("2006-01-02 15:04:05") + " " + i118Utils.I118Prt.Sprintf("found_scripts",
		color.CyanString(strconv.Itoa(total))))

	width := len(strconv.Itoa(len(arrs)))
	for idx, arr := range arrs {
		numb := fmt.Sprintf("#%0"+strconv.Itoa(width)+"d", idx+1)

		logUtils.PrintTo(logUtils.GetWholeLine(numb+" "+arr[3], "="))
		logUtils.PrintToWithColor(fmt.Sprintf("%s. %s", arr[0], arr[1]), color.FgCyan)

		fmt.Printf("Steps: \n%s \n", arr[2])

		logUtils.PrintToWithColor("", -1)
	}
}

func Brief(file string, keywords string) (bool, []string) {
	content := fileUtils.ReadFile(file)
	lang := langUtils.GetLangByFile(file)
	isOldFormat := strings.Index(content, "[esac]") > -1

	regStr := ""
	if isOldFormat {
		regStr = `\[case\][\S\s]*` +
			`title=([^\n]*)\n+` +
			`cid=([^\n]*)\n+` +
			`pid=([^\n]*)\n+` +
			`([\S\s]*)\n*` +
			`\[esac\]`
	} else {
		regStr = fmt.Sprintf(`(?sm)%s[\S\s]*`+
			`title=([^\n]*)\n+`+
			`cid=([^\n]*)\n+`+
			`pid=([^\n]*)\n+`+
			`([\S\s]*)\n*%s`,
			constant.LangCommentsRegxMap[lang][0], constant.LangCommentsRegxMap[lang][1])
	}
	myExp := regexp.MustCompile(regStr)
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
			return true, []string{caseId, title, steps, file}
		} else {
			return false, nil
		}
	}

	return false, nil
}

package action

import (
	"fmt"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/fatih/color"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func View(files []string, keywords string) {
	var cases []string
	for _, v1 := range files {
		group := scriptUtils.LoadScriptByWorkspace(v1)
		for _, v2 := range group {
			cases = append(cases, v2)
		}
	}

	view(cases, keywords)
}

func view(cases []string, keywords string) {
	keywords = strings.TrimSpace(keywords)
	count := 0

	arrs := make([][]string, 0)
	for _, file := range cases {
		pass, arr := brief(file, keywords)

		if pass {
			arrs = append(arrs, arr)
			count++
		}
	}

	total := len(arrs)

	logUtils.Info("\n" + time.Now().Format("2006-01-02 15:04:05") + " " + i118Utils.Sprintf("found_scripts",
		color.CyanString(strconv.Itoa(total))))

	width := len(strconv.Itoa(len(arrs)))
	for idx, arr := range arrs {
		numb := fmt.Sprintf("#%0"+strconv.Itoa(width)+"d", idx+1)

		logUtils.Infof(logUtils.GetWholeLine(numb+" "+arr[3], "="))
		logUtils.ExecConsole(color.FgCyan, fmt.Sprintf("%s. %s", arr[0], arr[1]))

		fmt.Printf("Steps: \n%s \n", arr[2])

		logUtils.Info("")
	}
}

func brief(file string, keywords string) (bool, []string) {
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
			langUtils.LangCommentsRegxMap[lang][0], langUtils.LangCommentsRegxMap[lang][1])
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

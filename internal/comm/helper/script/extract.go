package scriptUtils

import (
	"fmt"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"regexp"
	"strings"
)

const (
	groupTag = "group:"
	stepTag  = "step:"

	funcRegex               = `(?U)e\(['"](.+)['"]\)`
	singleLineCommentsRegex = `.*(?://|#)(.+)$`
	multiLineCommentsRegex  = `/\*+(.+)\*+/`
)

func Extract(scriptPaths []string) error {
	if len(scriptPaths) < 1 {
		logUtils.Infof("\n" + i118Utils.Sprintf("no_cases"))
		return nil
	}

	for _, pth := range scriptPaths {
		stepObjs := extractFromComments(pth)
		steps := prepareSteps(stepObjs)
		desc := prepareDesc(steps, pth)

		replaceCaseDesc(desc, pth)
	}

	return nil
}
func prepareSteps(stepObjs []*commDomain.ZtfStep) (steps []string) {
	for index, stepObj := range stepObjs {
		line := stepObj.Desc

		if len(stepObj.Children) == 0 && stepObj.Expect != "" {
			if stepObj.MultiLine {
				line += " >>\n" + stepObj.Expect + ">>\n"
			} else {
				line += " >> " + stepObj.Expect
			}
		}

		steps = append(steps, line)

		for _, childObj := range stepObj.Children {
			lineChild := "  " + childObj.Desc + " >> " + childObj.Expect
			steps = append(steps, lineChild)
		}

		if (index < len(stepObjs)-1 && len(stepObj.Children) == 0 && len(stepObjs[index+1].Children) > 0) || len(stepObj.Children) > 0 {
			steps = append(steps, "")
		}
	}

	return
}

func extractFromComments(file string) (stepObjs []*commDomain.ZtfStep) {
	lang := langUtils.GetLangByFile(file)
	content := fileUtils.ReadFile(file)

	findCaseTag := false
	start := false
	inGroup := false

	lines := strings.Split(content, "\n")
	for index := 0; index < len(lines); index++ {
		line := lines[index]
		lowerLine := strings.ToLower(line)

		if strings.Index(lowerLine, "cid") > -1 {
			findCaseTag = true
			continue
		}

		if findCaseTag && !start {
			if !isCommentEndTag(line, lang) {
				continue
			}

			start = true
			continue
		}
		if !start {
			continue
		}

		if strings.Index(lowerLine, groupTag) > -1 { // is group
			groupName, _ := getName(line, groupTag, lang)
			stepObj := commDomain.ZtfStep{Desc: groupName}
			stepObjs = append(stepObjs, &stepObj)
			inGroup = true

			continue
		}

		if strings.Index(lowerLine, stepTag) > -1 {
			isMuti := false
			stepName, expect := getName(line, stepTag, lang)
			if expect == "" {
				moreExpect, increase := parseMutiStep(lang, lines[index+1:])
				if increase > 0 {
					expect = moreExpect
					index += increase
					isMuti = true
				}
			}

			stepObj := commDomain.ZtfStep{Desc: stepName, Expect: expect, MultiLine: isMuti}

			if inGroup {
				if strings.Index(line, "]]") > -1 {
					inGroup = false
					stepObj.Expect = strings.TrimRight(strings.TrimSpace(stepObj.Expect), "]]")
				}

				stepObjs[len(stepObjs)-1].Children = append(stepObjs[len(stepObjs)-1].Children, stepObj)
			} else {
				stepObjs = append(stepObjs, &stepObj)
			}
		}

		// find e() function and its comments, for zentao user only
		myExp := regexp.MustCompile(funcRegex)
		arr := myExp.FindStringSubmatch(line)

		if len(arr) > 1 {
			expect := arr[1]
			desc := ""

			myExp := regexp.MustCompile(singleLineCommentsRegex)
			arr2 := myExp.FindStringSubmatch(line)
			if len(arr2) > 1 { // find single line comments on right
				desc = strings.TrimSpace(arr2[1])
			} else {
				preLine := strings.TrimSpace(lines[index-1])
				arr3 := myExp.FindStringSubmatch(preLine)
				if len(arr3) > 1 { // find single line comments on top
					desc = strings.TrimSpace(arr3[1])
				} else {
					myExp := regexp.MustCompile(multiLineCommentsRegex)
					arr4 := myExp.FindStringSubmatch(preLine)
					if len(arr4) > 1 { // find muti line comments on top
						desc = strings.TrimSpace(arr4[1])
					}
				}
			}

			stepObj := commDomain.ZtfStep{Desc: desc, Expect: expect, MultiLine: false}
			stepObjs = append(stepObjs, &stepObj)
		}
	}
	return
}

func prepareDesc(steps []string, file string) (desc string) {
	pass, caseId, productId, title := GetCaseInfo(file)
	if !pass {
		return
	}

	info := make([]string, 0)
	info = append(info, fmt.Sprintf("title=%s", title))
	info = append(info, fmt.Sprintf("cid=%d", caseId))
	info = append(info, fmt.Sprintf("pid=%d", productId))
	info = append(info, "\n"+strings.Join(steps, "\n"))

	desc = strings.Join(info, "\n")
	return
}

func replaceCaseDesc(desc, file string) {
	content := fileUtils.ReadFile(file)
	lang := langUtils.GetLangByFile(file)

	regStr := fmt.Sprintf(`(?smU)%s((?U:.*pid.*))\n(.*)%s`,
		langUtils.LangCommentsRegxMap[lang][0], langUtils.LangCommentsRegxMap[lang][1])

	re, _ := regexp.Compile(regStr)
	out := re.ReplaceAllString(content, "\n/**\n\n"+desc+"\n\n*/")

	fileUtils.WriteFile(file, out)
}

func getName(line, str, lang string) (name, expect string) {
	lowerLine := strings.ToLower(line)

	idx := strings.Index(lowerLine, str)
	name = line[idx+len(str):]

	if strings.Index(str, "step:") > -1 {
		arr := strings.Split(name, ">>")
		if len(arr) > 1 {
			name = strings.TrimSpace(arr[0])
			expect = strings.TrimSpace(arr[1])
		}
	}

	regx, _ := regexp.Compile(langUtils.LangCommentsTagMap[lang][1])
	name = strings.TrimSpace(regx.ReplaceAllString(name, ""))
	expect = strings.TrimSpace(regx.ReplaceAllString(expect, ""))

	return
}

func parseMutiStep(lang string, nextLines []string) (ret string, increase int) {
	for index, line := range nextLines {
		if isCommentStartTag(line, lang) || isCommentEndTag(line, lang) {
			break
		}

		if strings.TrimSpace(line) == ">>" {
			increase = index
			break
		}

		ret += "  " + strings.TrimSpace(line) + "\n"
	}

	if increase == 0 { // multi-line
		ret = ""
	}

	return
}

func isCommentStartTag(str, lang string) (pass bool) {
	pass, _ = regexp.MatchString(langUtils.LangCommentsRegxMap[lang][0], str)
	return
}

func isCommentEndTag(str, lang string) (pass bool) {
	pass, _ = regexp.MatchString(langUtils.LangCommentsRegxMap[lang][1], str)
	return
}

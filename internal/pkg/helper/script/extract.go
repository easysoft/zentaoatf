package scriptHelper

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

const (
	groupTag = "group:"
	stepTag  = "step:"

	funcRegex               = `(?U)\We\(['"](.+)['"]\)`
	singleLineCommentsRegex = `.*(?://|#)(.+)$`
	multiLineCommentsRegex  = `/\*+(.+)\*+/`
)

func Extract(scriptPaths []string) (done bool, err error) {
	if len(scriptPaths) < 1 {
		logUtils.Infof("\n" + i118Utils.Sprintf("no_cases"))
		return
	}

	for _, pth := range scriptPaths {
		stepObjs := extractFromComments(pth)
		steps := prepareSteps(stepObjs)
		steps = genDescFromRPE(pth)
		desc := prepareDesc(steps, pth)

		if steps != nil && len(steps) > 0 {
			ReplaceCaseDesc(desc, pth)
			done = true
		}
	}

	return
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
	lang := langHelper.GetLangByFile(file)
	content := fileUtils.ReadFile(file)

	findCaseTag := false
	start := false
	inGroup := false
	rpeDescLines := make([]string, 0)
	rpeIndex := 0

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
				moreExpect, increase := parseMultiStep(lang, lines[index+1:])
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
			//find function file and extract
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

			if desc == "" {
				if len(rpeDescLines) == 0 {
					rpeDescLines = genDescFromRPE(file)
				}
				if len(rpeDescLines) > rpeIndex {
					desc = rpeDescLines[rpeIndex]
				}
			}

			rpeIndex++
			stepObj := commDomain.ZtfStep{Desc: desc, Expect: expect, MultiLine: false}
			stepObjs = append(stepObjs, &stepObj)
		}
	}
	return
}

func prepareDesc(steps []string, file string) (desc string) {
	_, caseId, _, title, timeout := GetCaseInfo(file)
	// if !pass {
	// 	return
	// }

	info := make([]string, 0)
	info = append(info, fmt.Sprintf("title=%s", title))
	info = append(info, fmt.Sprintf("timeout=%d", timeout))
	info = append(info, fmt.Sprintf("cid=%d", caseId))
	//info = append(info, fmt.Sprintf("pid=%d", productId))
	info = append(info, "\n"+strings.Join(steps, "\n"))

	desc = strings.Join(info, "\n")
	return
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

	name = strings.TrimSpace(strings.Replace(name, commConsts.LangCommentsTagMap[lang][1], "", -1))
	expect = strings.TrimSpace(strings.Replace(expect, commConsts.LangCommentsTagMap[lang][1], "", -1))

	return
}

func parseMultiStep(lang string, nextLines []string) (ret string, increase int) {
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
	pass, _ = regexp.MatchString(commConsts.LangCommentsRegxMap[lang][0], str)
	return
}

func isCommentEndTag(str, lang string) (pass bool) {
	pass, _ = regexp.MatchString(commConsts.LangCommentsRegxMap[lang][1], str)
	return
}

func genDescFromRPE(file string) []string {
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/bash", "-c", file+" -extract")

	output := make([]string, 0)

	content, err := cmd.CombinedOutput()
	if err != nil {
		return output
	}

	output = strings.Split(string(content), "\n")

	return output
}

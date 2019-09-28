package scriptUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func Sort(cases []string) {
	for _, file := range cases {
		SortFile(file, false)
	}

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("success_sort_steps", len(cases)), -1)
}

func SortFile(script string, output bool) (map[string]string, map[string]string, map[string]string) {
	stepMap := make(map[string]string, 0)
	stepTypeMap := make(map[string]string, 0)
	expectMap := make(map[string]string, 0)

	if fileUtils.FileExist(script) {
		script := fileUtils.ReadFile(script)

		info, content := zentaoUtils.ReadCaseInfo(script)

		lines := strings.Split(content, "\n")

		groupBlockArr := getGroupBlockArr(lines)
		groupArr := getStepNestedArr(groupBlockArr)
		stepsTxt, stepMap, stepTypeMap, expectMap := getSortedTextFromNestedSteps(groupArr)

		// replace info
		re, _ := regexp.Compile(`(?s)\[case\].*\[esac\]`)
		script = re.ReplaceAllString(script, "[case]\n"+info+"\n"+stepsTxt+"\n\n[esac]")

		fileUtils.WriteFile(script, script)

		pass, expectIndependentContent := zentaoUtils.GetDependentExpect(script)
		if pass {
			expectMap = getExpectMapFromIndependentFile(stepMap, expectIndependentContent)
		}

	}

	return stepMap, stepTypeMap, expectMap
}

func getGroupBlockArr(lines []string) [][]string {
	groupBlockArr := make([][]string, 0)

	idx := 0
	for true {
		if idx >= len(lines) {
			break
		}

		var groupContent []string
		line := strings.TrimSpace(lines[idx])
		if isGroup(line) { // must match a group
			groupContent = make([]string, 0)
			groupContent = append(groupContent, line)

			idx++

			for true {
				if idx >= len(lines) {
					groupBlockArr = append(groupBlockArr, groupContent)
					break
				}

				line = strings.TrimSpace(lines[idx])
				if isGroup(line) {
					groupBlockArr = append(groupBlockArr, groupContent)

					break
				} else if line != "" && !isGroup(line) {
					groupContent = append(groupContent, line)
				}

				idx++
			}
		} else {
			idx++
		}
	}

	return groupBlockArr
}

func getStepNestedArr(blocks [][]string) []model.TestStep {
	ret := make([]model.TestStep, 0)

	for _, block := range blocks {
		name := block[0]
		group := model.TestStep{Desc: name}

		if isStepsIdent(block[1]) { // muti line
			group.MutiLine = true
			childs := loadMutiLineSteps(block[1:])

			group.Children = append(group.Children, childs...)
		} else {
			childs := loadSingleLineSteps(block[1:])

			group.Children = append(group.Children, childs...)
		}

		ret = append(ret, group)
	}

	return ret
}

func loadMutiLineSteps(arr []string) []model.TestStep {
	childs := make([]model.TestStep, 0)

	child := model.TestStep{}
	idx := 0
	for true {
		if idx >= len(arr) {
			if child.Desc != "" {
				childs = append(childs, child)
			}

			break
		}

		line := arr[idx]
		line = strings.TrimSpace(line)

		if isStepsIdent(line) {
			if idx > 0 {
				childs = append(childs, child)
			}

			child = model.TestStep{}
			idx++

			stp := ""
			for true {
				if idx >= len(arr) || hasBrackets(arr[idx]) {
					child.Desc = stp
					break
				}

				stp += arr[idx] + "\n"
				idx++
			}
		}

		if isExpectsIdent(line) {
			idx++

			exp := ""
			for true {
				if idx >= len(arr) || hasBrackets(arr[idx]) {
					child.Expect = exp
					break
				}

				exp += arr[idx] + "\n"
				idx++
			}
		}

	}

	return childs
}

func loadSingleLineSteps(arr []string) []model.TestStep {
	childs := make([]model.TestStep, 0)

	for _, line := range arr {
		line = strings.TrimSpace(line)

		sections := strings.Split(line, ">>")

		expect := ""
		if len(sections) > 1 {
			expect = sections[1]
		}

		child := model.TestStep{Desc: sections[0], Expect: expect}

		childs = append(childs, child)
	}

	return childs
}

func isGroupIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[\s*group\s*\]`, str)
	return pass
}

func isStepsIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[.*steps\.*\]`, str)
	return pass
}

func isExpectsIdent(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[.*expects\.*\]`, str)
	return pass
}

func hasBrackets(str string) bool {
	pass, _ := regexp.MatchString(`(?i)\[.*\]`, str)
	return pass
}

func isGroup(str string) bool {
	ret := hasBrackets(str) && !isStepsIdent(str) && !isExpectsIdent(str)

	return ret
}

func getSortedTextFromNestedSteps(groups []model.TestStep) (string, map[string]string, map[string]string, map[string]string) {
	ret := make([]string, 0)
	stepMap := make(map[string]string, 0)
	stepTypeMap := make(map[string]string, 0)
	expectMap := make(map[string]string, 0)

	groupNumb := 1
	for _, group := range groups {
		desc := group.Desc

		if desc == "[group]" {
			ret = append(ret, "\n"+desc)

			for idx, child := range group.Children {
				if group.MutiLine {
					// steps
					desc = replaceNumb("[steps]", groupNumb, -1, true)
					ret = append(ret, "  "+desc)

					ret = append(ret, printMutiStepOrExpect(child.Desc))

					// expects
					desc = replaceNumb("[expects]", groupNumb, -1, true)
					ret = append(ret, "  "+desc)

					ret = append(ret, printMutiStepOrExpect(child.Expect))
					if idx < len(group.Children)-1 {
						ret = append(ret, "")
					}
				} else {
					desc = replaceNumb(child.Desc, groupNumb, -1, false)
					desc = strings.TrimSpace(desc)

					expect := child.Expect
					expect = strings.TrimSpace(expect)
					if expect != "" {
						expect = ">> " + expect
					}

					ret = append(ret, fmt.Sprintf("  %s %s", desc, expect))
				}

				groupNumb++
			}
		} else {
			desc = replaceNumb(group.Desc, groupNumb, -1, true)
			ret = append(ret, "\n"+desc)

			childNumb := 1
			for _, child := range group.Children {
				if group.MutiLine {
					// steps
					desc = replaceNumb("[steps]", groupNumb, childNumb, true)
					ret = append(ret, "  "+desc)

					ret = append(ret, printMutiStepOrExpect(child.Desc))

					// expects
					desc = replaceNumb("[expects]", groupNumb, childNumb, true)
					ret = append(ret, "  "+desc)

					ret = append(ret, printMutiStepOrExpect(child.Expect))

					//if idx < len(group.Children)-1 {
					//	ret = append(ret, "")
					//}
				} else {
					desc = replaceNumb(child.Desc, groupNumb, -1, false)
					desc = strings.TrimSpace(desc)

					expect := child.Expect
					expect = strings.TrimSpace(expect)

					if expect != "" {
						expect = ">> " + expect
					}

					ret = append(ret, fmt.Sprintf("  %s %s", desc, expect))
				}

				childNumb++
			}

			groupNumb++
		}
	}

	return strings.Join(ret, "\n"), stepMap, stepTypeMap, expectMap
}

func replaceNumb(str string, groupNumb int, childNumb int, withBrackets bool) string {
	numb := getNumbStr(groupNumb, childNumb)

	reg := `[\d\.\s]*(.*)`
	repl := numb + " ${1}"
	if withBrackets {
		reg = `\[` + reg + `\]`
		repl = `[` + repl + `]`
	}

	regx, _ := regexp.Compile(reg)
	str = regx.ReplaceAllString(str, repl)

	return str
}
func getNumbStr(groupNumb int, childNumb int) string {
	numb := strconv.Itoa(groupNumb) + "."
	if childNumb != -1 {
		numb += strconv.Itoa(childNumb) + "."
	}

	return numb
}

func printMutiStepOrExpect(str string) string {
	str = strings.TrimSpace(str)

	ret := make([]string, 0)

	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)

		ret = append(ret, fmt.Sprintf("%s%s", strings.Repeat(" ", 4), line))
	}

	return strings.Join(ret, "\n")
}

func getExpectMapFromIndependentFile(stepMap map[string]string, content string) map[string]string {
	retMap := make(map[string]string, 0)

	expectArr := zentaoUtils.ReadExpectIndependentArr(content)

	idx := 0
	for k, _ := range stepMap {
		retMap[k] = strings.Join(expectArr[idx], "\r\n")
	}

	return retMap
}

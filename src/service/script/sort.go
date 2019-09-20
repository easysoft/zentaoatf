package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"regexp"
	"strconv"
	"strings"
)

func Sort(cases []string) {
	for _, file := range cases {
		if fileUtils.FileExist(file) {
			script := fileUtils.ReadFile(file)

			regx := regexp.MustCompile(`(?s)\[case\](?U)(.*)(\[.*)\[esac\]`)
			arr := regx.FindStringSubmatch(script)

			info := ""
			content := ""
			if len(arr) > 2 {
				info = strings.TrimSpace(arr[1])
				checkpoints := arr[2]
				content = commonUtils.RemoveBlankLine(checkpoints)
			}

			lines := strings.Split(content, "\n")

			groupBlock := getGroupBlock(lines)
			groupArr := getStepNestedArr(groupBlock)
			stepsTxt := getOrderTextFromNestedSteps(groupArr)

			// replace info
			re, _ := regexp.Compile(`(?s)\[case\].*\[esac\]`)
			script = re.ReplaceAllString(script, "[case]"+info+"\n"+stepsTxt+"\n\n[esac]")

			println(script)

			fileUtils.WriteFile(file, script)
			return
		}
	}

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("success_sort_steps", len(cases)), -1)
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

func getGroupBlock(lines []string) [][]string {
	groupBlockArr := make([][]string, 0)

	idx := 0
	for true {
		if idx >= len(lines) {
			break
		}

		var groupContent []string
		line := strings.TrimSpace(lines[idx])
		if isGroup(line) {
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
		}
	}

	return groupBlockArr
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

func getOrderTextFromNestedSteps(groups []model.TestStep) string {
	ret := make([]string, 0)

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
					ret = append(ret, fmt.Sprintf("  %s >> %s", desc, child.Expect))
				}

				groupNumb++
			}
		} else {
			desc = replaceNumb(group.Desc, groupNumb, -1, true)
			ret = append(ret, "\n"+desc)

			childNumb := 1
			for idx, child := range group.Children {
				if group.MutiLine {
					// steps
					desc = replaceNumb("[steps]", groupNumb, childNumb, true)
					ret = append(ret, "  "+desc)

					ret = append(ret, printMutiStepOrExpect(child.Desc))

					// expects
					desc = replaceNumb("[expects]", groupNumb, childNumb, true)
					ret = append(ret, "  "+desc)

					ret = append(ret, printMutiStepOrExpect(child.Expect))

					if idx < len(group.Children)-1 {
						ret = append(ret, "")
					}
				} else {
					desc = replaceNumb(child.Desc, groupNumb, -1, false)
					ret = append(ret, fmt.Sprintf("  %s >> %s", desc, child.Expect))
				}

				childNumb++
			}

			groupNumb++
		}
	}

	return strings.Join(ret, "\n")
}

func replaceNumb(str string, groupNumb int, childNumb int, withBrackets bool) string {
	numb := strconv.Itoa(groupNumb) + "."
	if childNumb != -1 {
		numb += strconv.Itoa(childNumb) + "."
	}

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

func printMutiStepOrExpect(str string) string {
	str = strings.TrimSpace(str)

	ret := make([]string, 0)

	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)

		ret = append(ret, fmt.Sprintf("%s%s", strings.Repeat(" ", 4), line))
	}

	return strings.Join(ret, "\n")
}

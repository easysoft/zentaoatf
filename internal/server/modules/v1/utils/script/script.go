package scriptUtils

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	resUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/res"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	"github.com/emirpasic/gods/maps"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GenerateScript(cs commDomain.ZtfCase, langType string, independentFile bool, caseIds *[]string,
	targetDir string, byModule bool, prefix string) {
	caseId := cs.Id
	productId := cs.Product
	moduleId := cs.Module
	caseTitle := cs.Title

	fileUtils.MkDirIfNeeded(targetDir)
	modulePath := ""
	if byModule && moduleId != "0" {
		modulePath = fmt.Sprintf("%d%s", moduleId, consts.PthSep)
	}

	content := ""
	isOldFormat := false
	scriptFile := fmt.Sprintf(targetDir+"%s%s%s.%s", modulePath, prefix, caseId, langUtils.LangMap[langType]["extName"])
	if fileUtils.FileExist(scriptFile) { // update title and steps
		content = fileUtils.ReadFile(scriptFile)
		isOldFormat = strings.Index(content, "[esac]") > -1
	}

	*caseIds = append(*caseIds, caseId)

	info := make([]string, 0)
	steps := make([]string, 0)
	independentExpects := make([]string, 0)
	srcCode := fmt.Sprintf("%s %s", langUtils.LangMap[langType]["commentsTag"],
		i118Utils.Sprintf("find_example", consts.PthSep, langType))

	info = append(info, fmt.Sprintf("title=%s", caseTitle))
	info = append(info, fmt.Sprintf("cid=%s", caseId))
	info = append(info, fmt.Sprintf("pid=%s", productId))

	StepWidth := 20
	stepDisplayMaxWidth := 0
	computerTestStepWidth(cs.StepArr, &stepDisplayMaxWidth, StepWidth)

	if isOldFormat {
		generateTestStepAndScriptObsolete(cs.StepArr, &steps, &independentExpects, independentFile)
	} else {
		generateTestStepAndScript(cs.StepArr, &steps, &independentExpects, independentFile)
	}
	info = append(info, strings.Join(steps, "\n"))

	if independentFile {
		expectFile := ScriptToExpectName(scriptFile)
		fileUtils.WriteFile(expectFile, strings.Join(independentExpects, "\n"))
	}

	if fileUtils.FileExist(scriptFile) { // update title and steps
		regStr := fmt.Sprintf(`(?sm)%s((?U:.*pid.*))\n(.*)%s`,
			langUtils.LangCommentsRegxMap[langType][0], langUtils.LangCommentsRegxMap[langType][1])

		// replace info
		re, _ := regexp.Compile(regStr)
		newContent := fmt.Sprintf("\n%s\n\n%s\n\n%s\n",
			langUtils.LangCommentsTagMap[langType][0],
			strings.Join(info, "\n"),
			langUtils.LangCommentsTagMap[langType][1])

		out := re.ReplaceAllString(content, newContent)

		fileUtils.WriteFile(scriptFile, out)
		return
	}

	path := fmt.Sprintf("res%stemplate%s", consts.PthSep, consts.PthSep)
	template, _ := resUtils.ReadRes(path + langType + ".tpl")

	out := fmt.Sprintf(string(template), strings.Join(info, "\n"), srcCode)
	fileUtils.WriteFile(scriptFile, out)
}

func generateTestStepAndScriptObsolete(testSteps []commDomain.ZtfStep, steps *[]string, independentExpects *[]string, independentFile bool) {
	nestedSteps := make([]commDomain.ZtfStep, 0)
	currGroup := commDomain.ZtfStep{}
	idx := 0

	// convert steps to nested
	for true {
		if idx >= len(testSteps) {
			break
		}

		ts := testSteps[idx]
		if ts.Parent == "0" && ts.Type != "group" { // flat step
			currGroup = commDomain.ZtfStep{Id: "-1", Desc: "group", Children: make([]commDomain.ZtfStep, 0)}
			currGroup.Children = append(currGroup.Children, ts)
			idx++

			mutiLine := false
			for true {
				if idx >= len(testSteps) {
					currGroup.MultiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}

				child := testSteps[idx]
				if child.Type != "group" { // flat step
					if !mutiLine {
						mutiLine = IsMultiLine(child)
					}

					currGroup.Children = append(currGroup.Children, child)
				} else { // found a group step
					currGroup.MultiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}
				idx++
			}
		} else if ts.Type == "group" {
			currGroup = commDomain.ZtfStep{Desc: ts.Desc, Children: make([]commDomain.ZtfStep, 0)}
			idx++

			mutiLine := false
			for true {
				if idx >= len(testSteps) {
					nestedSteps = append(nestedSteps, currGroup)
					break
				}

				child := testSteps[idx]
				if child.Type != "group" && child.Parent == ts.Id { // child step
					if !mutiLine {
						mutiLine = IsMultiLine(child)
					}

					currGroup.Children = append(currGroup.Children, child)
				} else { // found a group step
					currGroup.MultiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}
				idx++
			}
		}
	}

	stepNumb := 1
	// print nested steps, only one level
	for _, group := range nestedSteps {
		if group.Id == "-1" { // [group]
			*steps = append(*steps, fmt.Sprintf("\n[group]"))

			for _, child := range group.Children {
				*steps = append(*steps,
					GetCaseContent(child, strconv.Itoa(stepNumb), independentFile, group.MultiLine)...)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, getExcepts(child.Expect))
				}

				stepNumb++
			}
		} else { // [1. title]
			*steps = append(*steps, "\n"+fmt.Sprintf("[%d. %s]", stepNumb, group.Desc))

			for childNo, child := range group.Children {
				numbStr := fmt.Sprintf("%d.%d", stepNumb, childNo+1)
				*steps = append(*steps, GetCaseContent(child, numbStr, independentFile, group.MultiLine)...)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, getExcepts(child.Expect))
				}
			}

			stepNumb++
		}
	}
}

func generateTestStepAndScript(testSteps []commDomain.ZtfStep, steps *[]string, independentExpects *[]string, independentFile bool) {
	nestedSteps := make([]commDomain.ZtfStep, 0)

	// convert steps to nested
	for index := 0; index < len(testSteps); index++ {
		ts := testSteps[index]
		item := commDomain.ZtfStep{Desc: ts.Desc, Expect: ts.Expect, Children: make([]commDomain.ZtfStep, 0)}

		if ts.Type == "group" {
			nestedSteps = append(nestedSteps, item)
		} else if ts.Type == "item" {
			nestedSteps[len(nestedSteps)-1].Children = append(nestedSteps[len(nestedSteps)-1].Children, item)
		} else if ts.Type == "step" {
			nestedSteps = append(nestedSteps, item)
		}
	}

	// print nested steps, only one level
	stepNumb := 1
	*steps = append(*steps, "")
	for _, item := range nestedSteps {
		numbStr := fmt.Sprintf("%d", stepNumb)
		*steps = append(*steps, GetCaseContent(item, numbStr, independentFile, false)...)

		for childNo, child := range item.Children {
			numbStr := fmt.Sprintf("%d.%d", stepNumb, childNo+1)
			*steps = append(*steps, GetCaseContent(child, numbStr, independentFile, true)...)

			if independentFile && strings.TrimSpace(child.Expect) != "" {
				*independentExpects = append(*independentExpects, getExcepts(child.Expect))
			}
		}

		stepNumb++
	}
}

func computerTestStepWidth(steps []commDomain.ZtfStep, stepSDisplayMaxWidth *int, stepWidth int) {
	for _, ts := range steps {
		length := len(ts.Id)
		if length > *stepSDisplayMaxWidth {
			*stepSDisplayMaxWidth = length
		}
	}
	*stepSDisplayMaxWidth += stepWidth // prefix space and @step
}

func GenSuite(cases []string, targetDir string) {
	str := strings.Join(cases, "\n")

	fileUtils.WriteFile(targetDir+"all."+commConsts.ExtNameSuite, str)
}

func getExcepts(str string) string {
	str = stringUtils.TrimAll(str)

	arr := strings.Split(str, "\n")

	if len(arr) == 1 {
		return ">> " + str
	} else {
		return ">>\n" + str
	}
}

func GetStepAndExpectMap(file string) (stepMap, stepTypeMap, expectMap maps.Map, isOldFormat bool) {
	if !fileUtils.FileExist(file) {
		return
	}

	lang := langUtils.GetLangByFile(file)
	txt := fileUtils.ReadFile(file)

	isOldFormat = strings.Index(txt, "[esac]") > -1
	_, checkpoints := ReadCaseInfo(txt, lang, isOldFormat)
	lines := strings.Split(checkpoints, "\n")

	if isOldFormat {
		groupBlockArr := getGroupBlockArr(lines)
		groupArr := getStepNestedArrObsolete(groupBlockArr)
		_, stepMap, stepTypeMap, expectMap = getSortedTextFromNestedStepsObsolete(groupArr)
	} else {
		groupArr := getStepNestedArr(lines)
		_, stepMap, stepTypeMap, expectMap = getSortedTextFromNestedSteps(groupArr)
	}

	return
}

func SortFile(file string) {
	stepsTxt := ""

	if fileUtils.FileExist(file) {
		txt := fileUtils.ReadFile(file)
		lang := langUtils.GetLangByFile(file)
		isOldFormat := strings.Index(txt, "[esac]") > -1
		info, content := ReadCaseInfo(txt, lang, isOldFormat)
		lines := strings.Split(content, "\n")

		groupBlockArr := getGroupBlockArr(lines)
		groupArr := getStepNestedArrObsolete(groupBlockArr)
		stepsTxt, _, _, _ = getSortedTextFromNestedStepsObsolete(groupArr)

		// replace info
		from := ""
		to := ""
		if isOldFormat {
			from = `(?s)\[case\].*\[esac\]`
			to = "[case]\n" + info + "\n" + stepsTxt + "\n\n[esac]"
		} else {
			from = fmt.Sprintf(`(?s)%s.*%s`, langUtils.LangCommentsRegxMap[lang][0], langUtils.LangCommentsRegxMap[lang][1])
			to = fmt.Sprintf("%s\n"+info+"\n"+stepsTxt+"\n\n%s",
				langUtils.LangCommentsRegxMap[lang][0], langUtils.LangCommentsRegxMap[lang][1])
		}
		re, _ := regexp.Compile(from)
		script := re.ReplaceAllString(txt, to)

		fileUtils.WriteFile(file, script)
	}
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

func getStepNestedArrObsolete(blocks [][]string) (ret []commDomain.ZtfStep) {
	for _, block := range blocks {
		name := block[0]
		group := commDomain.ZtfStep{Desc: name}

		if isStepsIdent(block[1]) { // muti line
			group.MultiLine = true
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

func getStepNestedArr(lines []string) (ret []commDomain.ZtfStep) {
	parent := commDomain.ZtfStep{}
	increase := 0
	for index := 0; index < len(lines); index++ {
		line := lines[index]
		lineTrim := strings.TrimSpace(line)
		if lineTrim == "" || lineTrim == ">>" {
			continue
		}

		if strings.Index(line, " ") != 0 {
			parent, increase = parserNextLines(line, lines[index+1:])
			index += increase

			if strings.TrimSpace(parent.Expect) == "" && strings.Index(line, ">>") > -1 {
				parent.Expect = commConsts.ExpectResultPass
			}
			ret = append(ret, parent)
		} else { // 有缩进
			child := commDomain.ZtfStep{}
			child, increase = parserNextLines(line, lines[index+1:])
			index += increase

			if parent.Desc != "" {
				if strings.TrimSpace(child.Expect) == "" && strings.Index(line, ">>") > -1 {
					child.Expect = commConsts.ExpectResultPass
				}

				ret[len(ret)-1].Children = append(ret[len(ret)-1].Children, child)
			}
		}
	}

	return
}
func parserNextLines(str string, nextLines []string) (ret commDomain.ZtfStep, increase int) {
	arr := strings.Split(str, ">>")
	desc := strings.TrimSpace(arr[0])

	expect := ""
	if len(arr) > 1 {
		expect = strings.TrimSpace(arr[1])
	}

	if strings.Index(str, ">>") < 0 || expect != "" { // no >> or single line expect
		ret = commDomain.ZtfStep{Desc: desc, Expect: expect}
		return
	}

	if strings.Index(str, ">>") > -1 { // will test if it has multi-line expect
		for index, line := range nextLines {
			if strings.TrimSpace(line) == ">>" {
				increase = index
				break
			}

			if strings.Index(line, ">>") > -1 {
				expect = ""
				break
			}

			if len(expect) > 0 {
				expect += " | "
			}
			expect += strings.TrimSpace(line)
		}

		if increase == 0 { // multi-line
			expect = ""
		}
	}

	ret = commDomain.ZtfStep{Desc: desc, Expect: expect}
	return
}

func loadMutiLineSteps(arr []string) []commDomain.ZtfStep {
	childs := make([]commDomain.ZtfStep, 0)

	child := commDomain.ZtfStep{}
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

			child = commDomain.ZtfStep{}
			idx++

			stp := ""
			for true { // retrieve next lines
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
			for true { // retrieve next lines
				if idx >= len(arr) || hasBrackets(arr[idx]) {
					child.Expect = exp
					break
				}

				temp := strings.TrimSpace(arr[idx])
				if temp == ">>" {
					temp = ""
				}
				exp += temp + "\n"
				idx++
			}
		}

	}

	return childs
}

func loadSingleLineSteps(arr []string) []commDomain.ZtfStep {
	childs := make([]commDomain.ZtfStep, 0)

	for _, line := range arr {
		line = strings.TrimSpace(line)

		sections := strings.Split(line, ">>")
		expect := ""
		if len(sections) > 1 { // has expect
			expect = strings.TrimSpace(sections[1])
		}

		child := commDomain.ZtfStep{Desc: sections[0], Expect: expect}

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
	pass, _ := regexp.MatchString(`(?i)()\[.*\]`, str)
	return pass
}

func isGroup(str string) bool {
	ret := strings.Index(str, ">>") < 0 && hasBrackets(str) && !isStepsIdent(str) && !isExpectsIdent(str)

	return ret
}

func getSortedTextFromNestedStepsObsolete(groups []commDomain.ZtfStep) (string, maps.Map, maps.Map, maps.Map) {
	ret := make([]string, 0)
	stepMap := linkedhashmap.New()
	stepTypeMap := linkedhashmap.New()
	expectMap := linkedhashmap.New()

	groupNumb := 1
	for _, group := range groups {
		desc := group.Desc

		if desc == "[group]" {
			ret = append(ret, "\n"+desc)

			for idx, child := range group.Children { // level 1 item
				numbStr := getNumbStr(groupNumb, -1)
				stepTypeMap.Put(numbStr, "item")

				if group.MultiLine {
					// steps
					tag := replaceNumb("[steps]", groupNumb, -1, true)
					ret = append(ret, "  "+tag)

					stepTxt := printMutiStepOrExpect(child.Desc)
					ret = append(ret, stepTxt)
					stepMap.Put(numbStr, stepTxt)

					// expects
					tag = replaceNumb("[expects]", groupNumb, -1, true)
					ret = append(ret, "  "+tag)

					expectTxt := printMutiStepOrExpect(child.Expect)
					ret = append(ret, expectTxt)
					if idx < len(group.Children)-1 {
						ret = append(ret, "")
					}
					expectMap.Put(numbStr, expectTxt)
				} else {
					stepTxt := strings.TrimSpace(child.Desc)
					stepTxtWithNumb := replaceNumb(stepTxt, groupNumb, -1, false)
					stepMap.Put(numbStr, stepTxt)

					expectTxt := child.Expect
					expectTxt = strings.TrimSpace(expectTxt)
					expectMap.Put(numbStr, expectTxt)

					if expectTxt != "" {
						expectTxt = ">> " + expectTxt
					}

					ret = append(ret, fmt.Sprintf("  %s %s", stepTxtWithNumb, expectTxt))
				}

				groupNumb++
			}
		} else {
			desc = replaceNumb(group.Desc, groupNumb, -1, true)
			ret = append(ret, "\n"+desc)

			numbStr := getNumbStr(groupNumb, -1)
			stepMap.Put(numbStr, getGroupName(group.Desc))
			stepTypeMap.Put(numbStr, "group")
			expectMap.Put(numbStr, "")

			childNumb := 1
			for _, child := range group.Children {
				numbStr := getNumbStr(groupNumb, childNumb)
				stepTypeMap.Put(numbStr, "item")

				if group.MultiLine {
					// steps
					tag := replaceNumb("[steps]", groupNumb, childNumb, true)
					ret = append(ret, "  "+tag)

					stepTxt := printMutiStepOrExpect(child.Desc)
					ret = append(ret, stepTxt)
					stepMap.Put(numbStr, stepTxt)

					// expects
					tag = replaceNumb("[expects]", groupNumb, childNumb, true)
					ret = append(ret, "  "+tag)

					expectTxt := printMutiStepOrExpect(child.Expect)
					ret = append(ret, expectTxt)
					expectMap.Put(numbStr, expectTxt)
				} else {
					stepTxt := strings.TrimSpace(child.Desc)
					stepMap.Put(numbStr, stepTxt)

					expectTxt := child.Expect
					expectTxt = strings.TrimSpace(expectTxt)
					expectMap.Put(numbStr, expectTxt)

					if expectTxt != "" {
						expectTxt = ">> " + expectTxt
					}

					ret = append(ret, fmt.Sprintf("  %s %s", stepTxt, expectTxt))
				}

				childNumb++
			}

			groupNumb++
		}
	}

	return strings.Join(ret, "\n"), stepMap, stepTypeMap, expectMap
}

func getSortedTextFromNestedSteps(groups []commDomain.ZtfStep) (ret string, stepMap, stepTypeMap, expectMap maps.Map) {
	arr := make([]string, 0)
	stepMap = linkedhashmap.New()
	stepTypeMap = linkedhashmap.New()
	expectMap = linkedhashmap.New()

	for idx1, group := range groups {
		numbStr := getNumbStr(idx1+1, -1)
		stepType := "step"
		if len(group.Children) > 0 {
			stepType = "group"
		}
		stepTypeMap.Put(numbStr, stepType)

		stepTxt := strings.TrimSpace(group.Desc)
		stepMap.Put(numbStr, stepTxt)

		expectTxt := strings.TrimSpace(group.Expect)
		expectTxt = strings.TrimRight(expectTxt, "]]")
		expectTxt = strings.TrimSpace(expectTxt)

		expectMap.Put(numbStr, expectTxt)

		if expectTxt != "" {
			expectTxt = ">> " + expectTxt
		}
		arr = append(arr, fmt.Sprintf("  %s %s", stepTxt, expectTxt))

		for idx2, child := range group.Children {
			numbStr := getNumbStr(idx1+1, idx2+1)
			stepTypeMap.Put(numbStr, "item")

			stepTxt := strings.TrimSpace(child.Desc)
			stepMap.Put(numbStr, stepTxt)

			expectTxt := strings.TrimSpace(child.Expect)
			expectMap.Put(numbStr, expectTxt)

			if expectTxt != "" {
				expectTxt = ">> " + expectTxt
			}

			arr = append(arr, fmt.Sprintf("  %s %s", stepTxt, expectTxt))
		}
	}

	ret = strings.Join(arr, "\n")
	return
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
func getGroupName(str string) string {
	reg := `\[\d\.\s]*(.*)\]`
	repl := "${1}"

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

	return strings.Join(ret, "\r\n")
}

func GetExpectMapFromIndependentFileObsolete(expectMap maps.Map, content string, withEmptyExpect bool) maps.Map {
	retMap := linkedhashmap.New()

	expectArr := ReadExpectIndependentArrObsolete(content)

	idx := 0
	for _, keyIfs := range expectMap.Keys() {
		valueIfs, _ := expectMap.Get(keyIfs)

		key := strings.TrimSpace(keyIfs.(string))
		value := strings.TrimSpace(valueIfs.(string))

		if value == "" && len(expectArr) > idx {
			retMap.Put(key, strings.Join(expectArr[idx], "\r\n"))
			idx++
		} else {
			if withEmptyExpect {
				retMap.Put(key, "")
			}
		}
	}

	return retMap
}

func GetExpectMapFromIndependentFile(expectMap maps.Map, content string, withEmptyExpect bool) maps.Map {
	retMap := linkedhashmap.New()

	expectArr := ReadExpectIndependentArr(content)

	idx := 0
	for _, keyIfs := range expectMap.Keys() {
		valueIfs, _ := expectMap.Get(keyIfs)

		key := strings.TrimSpace(keyIfs.(string))
		value := strings.TrimSpace(valueIfs.(string))

		if value == "" && len(expectArr) > idx {
			retMap.Put(key, strings.Join(expectArr[idx], " | "))
			idx++
		} else {
			if withEmptyExpect {
				retMap.Put(key, "")
			}
		}
	}

	return retMap
}

func GetCaseContent(stepObj commDomain.ZtfStep, seq string, independentFile bool, isChild bool) []string {
	lines := make([]string, 0)

	step := strings.TrimSpace(stepObj.Desc)
	expect := strings.TrimSpace(stepObj.Expect)

	stepStr := getStepContent(step, isChild)
	expectStr := getExpectContent(expect, isChild, independentFile)

	lines = append(lines, stepStr+expectStr)

	return lines
}

func getStepContent(str string, isChild bool) (ret string) {
	str = strings.TrimSpace(str)

	rpl := "\n"
	if isChild {
		rpl = "\n" + "  "
	}
	ret = strings.ReplaceAll(str, "\r\n", rpl)
	if isChild {
		ret = "  " + ret
	}

	return
}
func getExpectContent(str string, isChild bool, independentFile bool) (ret string) {
	str = strings.TrimSpace(str)
	if str == "" {
		return
	}

	isMultiLine := strings.Count(str, "\r\n") > 0
	if !isMultiLine {
		if independentFile {
			ret = str
		} else {
			ret = " >> " + str
		}
	} else {
		rpl := "\r\n" + "  "

		if independentFile {
			ret = ">>\n" + strings.ReplaceAll(str, "\r\n", rpl) + "\n>>"
		} else {
			ret = " >> " + strings.ReplaceAll(str, "\r\n", rpl) + "\n>>"
		}
	}

	return
}

func IsMultiLine(step commDomain.ZtfStep) bool {
	if strings.Index(step.Desc, "\n") > -1 || strings.Index(step.Expect, "\n") > -1 {
		return true
	}

	return false
}

func ScriptToExpectName(file string) string {
	fileSuffix := path.Ext(file)
	expectName := strings.TrimSuffix(file, fileSuffix) + ".exp"

	return expectName
}

//func RunDateFolder() string {
//	runName := dateUtils.DateTimeStrFmt(time.Now(), "2006-01-02T150405") + string(os.PathSeparator)
//
//	return runName
//}

func GetCaseInfo(file string) (bool, int, int, string) {
	var caseId int
	var productId int
	var title string

	content := fileUtils.ReadFile(file)
	isOldFormat := strings.Index(content, "[esac]") > -1
	pass := CheckFileContentIsScript(content)
	if !pass {
		return false, caseId, productId, title
	}

	caseInfo := ""
	lang := langUtils.GetLangByFile(file)
	regStr := ""
	if isOldFormat {
		regStr = `(?s)\[case\](.*)\[esac\]`
	} else {
		regStr = fmt.Sprintf(`(?sm)%s((?U:.*pid.*))\n(.*)%s`,
			langUtils.LangCommentsRegxMap[lang][0], langUtils.LangCommentsRegxMap[lang][1])
	}
	myExp := regexp.MustCompile(regStr)
	arr := myExp.FindStringSubmatch(content)
	if len(arr) > 1 {
		caseInfo = arr[1]
	}

	caseInfo += "\n"

	myExp = regexp.MustCompile(`[\S\s]*cid=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		caseId, _ = strconv.Atoi(arr[1])
	}

	myExp = regexp.MustCompile(`[\S\s]*pid=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		productId, _ = strconv.Atoi(arr[1])
	}

	myExp = regexp.MustCompile(`[\S\s]*title=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		title = arr[1]
	}

	return pass, caseId, productId, title
}

//func ReadScriptCheckpoints(file string) ([]string, [][]string) {
//	_, expectIndependentContent := GetDependentExpect(file)
//
//	content := fileUtils.ReadFile(file)
//	_, checkpoints := ReadCaseInfo(content)
//
//	cpStepArr, expectArr := getCheckpointStepArr(checkpoints, expectIndependentContent)
//
//	return cpStepArr, expectArr
//}
//func getCheckpointStepArr(content string, expectIndependentContent string) ([]string, [][]string) {
//	cpStepArr := make([]string, 0)
//	expectArr := make([][]string, 0)
//
//	independentExpect := expectIndependentContent != ""
//
//	lines := strings.Split(content, "\n")
//	i := 0
//	for i < len(lines) {
//		step := ""
//		expects := make([]string, 0)
//
//		line := strings.TrimSpace(lines[i])
//
//		regx := regexp.MustCompile(`(?U:[\d\.]*)(.+)>>(.*)`)
//		arr := regx.FindStringSubmatch(line)
//		if len(arr) > 2 {
//			step = arr[1]
//			if !independentExpect {
//				expects = append(expects, strings.TrimSpace(arr[2]))
//			}
//		} else {
//			regx = regexp.MustCompile(`\[([\d\.]*).*expects\]`)
//			arr = regx.FindStringSubmatch(line)
//			if len(arr) > 1 {
//				step = arr[1]
//
//				if !independentExpect {
//					for i+1 < len(lines) {
//						ln := strings.TrimSpace(lines[i+1])
//
//						if strings.Index(ln, "[") == 0 || strings.Index(ln, ">>") > 0 || ln == "" {
//							break
//						} else {
//							expects = append(expects, ln)
//							i++
//						}
//					}
//				}
//			}
//		}
//
//		if step != "" && len(expects) > 0 {
//			cpStepArr = append(cpStepArr, step)
//			if !independentExpect {
//				expectArr = append(expectArr, expects)
//			}
//		}
//		i++
//	}
//
//	if independentExpect {
//		expectArr = ReadExpectIndependentArrObsolete(expectIndependentContent)
//	}
//
//	return cpStepArr, expectArr
//}

func ReadExpectIndependentArrObsolete(content string) [][]string {
	lines := strings.Split(content, "\n")

	ret := make([][]string, 0)
	var cpArr []string

	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == ">>" { // more than one line
			cpArr = make([]string, 0)
		} else if strings.Index(line, ">>") == 0 { // single line
			line = strings.Replace(line, ">>", "", -1)
			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		} else { // under >>
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], ">>") > -1 {
				ret = append(ret, cpArr)
				cpArr = make([]string, 0)
			}
		}
	}

	return ret
}

func ReadExpectIndependentArr(content string) [][]string {
	lines := strings.Split(content, "\n")

	ret := make([][]string, 0)
	var cpArr []string

	model := ""
	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == ">>" { // more than one line
			model = "multi"
			cpArr = make([]string, 0)
		} else if model == "multi" { // in >> and >> in multi line mode
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], ">>") > -1 {
				temp := make([]string, 0)
				temp = append(temp, strings.Join(cpArr, " | "))

				ret = append(ret, temp)
				cpArr = make([]string, 0)
				model = ""
			}
		} else if line == ">>" {
			continue
		} else {
			model = "single"

			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		}
	}

	return ret
}

func ReadLogArrObsolete(content string) (isSkip bool, ret [][]string) {
	lines := strings.Split(content, "\n")

	ret = make([][]string, 0)
	var cpArr []string

	model := ""
	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line == "skip" {
			isSkip = true
			return
		}

		if line == ">>" { // more than one line
			model = "multi"
			cpArr = make([]string, 0)
		} else if strings.Index(line, ">>") == 0 { // single line
			model = "single"

			line = strings.Replace(line, ">>", "", -1)
			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		} else {
			if model == "" || model == "single" {
				continue
			}

			// under >>
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], ">>") > -1 {
				ret = append(ret, cpArr)
				cpArr = make([]string, 0)
			}
		}
	}

	return
}

func ReadLogArr(content string) (isSkip bool, ret [][]string) {
	lines := strings.Split(content, "\n")

	ret = make([][]string, 0)
	var cpArr []string

	model := ""
	for idx := 0; idx < len(lines); idx++ {
		line := strings.TrimSpace(lines[idx])

		if line == "skip" {
			isSkip = true
			return
		}

		if line == ">>" { // more than one line
			model = "multi"
			cpArr = make([]string, 0)
		} else if model == "multi" { // in >> and >> in multi line mode
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], ">>") > -1 {
				temp := make([]string, 0)
				temp = append(temp, strings.Join(cpArr, " | "))

				ret = append(ret, temp)
				cpArr = make([]string, 0)

				idx = idx + 1
				model = ""
			}
		} else if line == ">>" {
			continue
		} else {
			model = "single"

			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		}
	}

	return
}

func CheckFileIsScript(path string) bool {
	content := fileUtils.ReadFile(path)

	pass := CheckFileContentIsScript(content)

	return pass
}

func CheckFileContentIsScript(content string) bool {
	pass, _ := regexp.MatchString(`cid\b*=`, content)

	return pass
}

func ReadCaseInfo(content, lang string, isOldFormat bool) (info, checkpoints string) {
	regStr := ""
	if isOldFormat {
		regStr = `(?s)\[case\]((?U:.*pid.*))\n(.*)\[esac\]`
	} else {
		regStr = fmt.Sprintf(`(?smU)%s((?U:.*pid.*))\n(.*)%s`,
			langUtils.LangCommentsRegxMap[lang][0], langUtils.LangCommentsRegxMap[lang][1])
	}
	myExp := regexp.MustCompile(regStr)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 2 {
		info = strings.TrimSpace(arr[1])
		checkpoints = strings.TrimSpace(arr[2])

		return
	}

	return
}
func ReadCaseId(content string) string {
	myExp := regexp.MustCompile(`(?s).*\ncid=((?U:.*))\n.*`)
	arr := myExp.FindStringSubmatch(content)

	if len(arr) > 1 {
		id := strings.TrimSpace(arr[1])
		return id
	}

	return ""
}

func GetDependentExpect(file string) (bool, string) {
	dir := fileUtils.AddPathSepIfNeeded(filepath.Dir(file))
	name := strings.Replace(filepath.Base(file), path.Ext(file), ".exp", -1)
	expectIndependentFile := dir + name

	if !fileUtils.FileExist(expectIndependentFile) {
		expectIndependentFile = dir + "." + name
	}

	if fileUtils.FileExist(expectIndependentFile) {
		expectIndependentContent := fileUtils.ReadFile(expectIndependentFile)
		return true, expectIndependentContent
	}

	return false, ""
}

func GetScriptByIdsInDir(dirPth string, idMap map[int]string, files *[]string) error {
	dirPth = fileUtils.AbsolutePath(dirPth)

	sep := string(os.PathSeparator)

	if commonUtils.IgnoreFile(dirPth) {
		return nil
	}

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			GetScriptByIdsInDir(dirPth+name+sep, idMap, files)
		} else {
			regx := langUtils.GetSupportLanguageExtRegx()
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", name)

			if !pass {
				continue
			}

			path := dirPth + name

			pass, id, _, _ := GetCaseInfo(path)
			if pass {
				_, ok := idMap[id]

				if ok {
					*files = append(*files, path)
				}
			}
		}
	}

	return nil
}

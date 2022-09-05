package scriptHelper

import (
	"fmt"
	"html"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	"github.com/easysoft/zentaoatf/pkg/consts"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
)

func ReplaceCaseDesc(desc, file string) {
	content := fileUtils.ReadFile(file)
	lang := langHelper.GetLangByFile(file)

	regStr := fmt.Sprintf(`(?smU)%s((?U:.*pid.*))\n(.*)%s`,
		commConsts.LangCommentsRegxMap[lang][0], commConsts.LangCommentsRegxMap[lang][1])
	re, _ := regexp.Compile(regStr)

	newDesc := fmt.Sprintf("\n%s\n\n"+desc+"\n\n%s",
		commConsts.LangCommentsTagMap[lang][0],
		commConsts.LangCommentsTagMap[lang][1])

	out := re.ReplaceAllString(content, newDesc)

	fileUtils.WriteFile(file, out)
}

func GetStepAndExpectMap(file string) (steps []commDomain.ZentaoCaseStep) {
	if !fileUtils.FileExist(file) {
		return
	}

	lang := langHelper.GetLangByFile(file)
	txt := fileUtils.ReadFile(file)

	_, checkpoints := ReadCaseInfo(txt, lang)
	lines := strings.Split(checkpoints, "\n")

	groupArr := getStepNestedArr(lines)
	_, steps = getSortedTextFromNestedSteps(groupArr)

	isIndependent, expectIndependentContent := GetDependentExpect(file)
	if isIndependent {
		GetExpectMapFromIndependentFile(&steps, expectIndependentContent, false)
	}

	return
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
				expect += "\r\n"
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

func loadMultiLineSteps(arr []string) []commDomain.ZtfStep {
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
	children := make([]commDomain.ZtfStep, 0)

	for _, line := range arr {
		line = strings.TrimSpace(line)

		sections := strings.Split(line, ">>")
		expect := ""
		if len(sections) > 1 { // has expect
			expect = strings.TrimSpace(sections[1])
		}

		child := commDomain.ZtfStep{Desc: sections[0], Expect: expect}

		children = append(children, child)
	}

	return children
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

func getSortedTextFromNestedSteps(groups []commDomain.ZtfStep) (ret string, steps []commDomain.ZentaoCaseStep) {
	arr := make([]string, 0)

	for _, group := range groups {
		step := commDomain.ZentaoCaseStep{}

		stepType := commConsts.Item
		if len(group.Children) > 0 {
			stepType = commConsts.Group
		}
		step.Type = stepType

		stepTxt := strings.TrimSpace(group.Desc)
		step.Desc = stepTxt

		expectTxt := strings.TrimSpace(group.Expect)
		expectTxt = strings.TrimRight(expectTxt, "]]")
		expectTxt = strings.TrimSpace(expectTxt)

		step.Expect = expectTxt

		steps = append(steps, step)

		if expectTxt != "" {
			expectTxt = ">> " + expectTxt
		}
		arr = append(arr, fmt.Sprintf("  %s %s", stepTxt, expectTxt))

		for _, child := range group.Children {
			stepChild := commDomain.ZentaoCaseStep{}

			stepChild.Type = commConsts.Item

			stepTxt := strings.TrimSpace(child.Desc)
			stepChild.Desc = stepTxt

			expectTxt := strings.TrimSpace(child.Expect)
			stepChild.Expect = expectTxt

			steps = append(steps, stepChild)

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

func GetExpectMapFromIndependentFile(steps *[]commDomain.ZentaoCaseStep, content string, withEmptyExpect bool) {
	expectArr := ReadExpectIndependentArr(content)

	index := 0
	for idx, _ := range *steps {
		if len(expectArr) > index && (*steps)[idx].Expect == "pass" { // not set step that has no expect
			(*steps)[idx].Expect = strings.Join(expectArr[index], "\r\n")
			index++
		} else {
			if withEmptyExpect {
				(*steps)[idx].Expect = ""
			}
		}
	}

	return
}

func GetCaseContent(stepObj commDomain.ZtfStep, seq string, independentFile bool, isChild bool) (
	stepContent, expectContent string) {

	step := strings.TrimSpace(stepObj.Desc)
	expect := strings.TrimSpace(stepObj.Expect)

	stepStr := getStepContent(step, isChild)
	expectStr := getExpectContent(expect, isChild, independentFile)

	if !independentFile {
		stepContent = stepStr + expectStr
	} else {
		stepContent = stepStr
		if stepObj.Children == nil || len(stepObj.Children) == 0 {
			stepContent += " >>"
		}
	}

	expectContent = expectStr

	stepContent = html.UnescapeString(stepContent)
	expectContent = html.UnescapeString(expectContent)

	return
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

	isSingleLine := strings.Count(str, "\r\n") == 0
	if isSingleLine {
		if independentFile {
			ret = str
		} else {
			ret = " >> " + str
		}
	} else { // multi-line
		rpl := "\r\n"

		space := "  "
		spaceBeforeTerminator := ""
		spaceBeforeText := space
		if isChild {
			spaceBeforeTerminator = space
			spaceBeforeText = strings.Repeat(space, 2)
		}

		if independentFile {
			//>>
			//	expect 1.2 line 1
			//	expect 1.2 line 2
			//>>
			ret = ">>\n" + space + strings.ReplaceAll(str, rpl, rpl+space) + "\n>>"
		} else {
			//step 1.2 >>
			//	expect 1.2 line 1
			//  expect 1.2 line 2
			//>>
			ret = " >> \n" + spaceBeforeText +
				strings.ReplaceAll(str, rpl, rpl+spaceBeforeText) +
				"\n" + spaceBeforeTerminator + ">>"
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

func GetCaseInfo(file string) (pass bool, caseId, productId int, title string, timeout int64) {
	content := fileUtils.ReadFile(file)
	isOldFormat := strings.Index(content, "[esac]") > -1
	pass = CheckFileContentIsScript(content)
	if !pass {
		return false, caseId, productId, title, timeout
	}

	caseInfo := ""
	lang := langHelper.GetLangByFile(file)
	regStr := ""
	if isOldFormat {
		regStr = `(?s)\[case\](.*)\[esac\]`
	} else {
		regStr = fmt.Sprintf(`(?sm)%s((?U:.*pid.*))\n(.*)%s`,
			commConsts.LangCommentsRegxMap[lang][0], commConsts.LangCommentsRegxMap[lang][1])
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

	myExp = regexp.MustCompile(`[\S\s]*timeout=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		timeout, _ = strconv.ParseInt(arr[1], 10, 64)
	}

	myExp = regexp.MustCompile(`[\S\s]*pid=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		productId, _ = strconv.Atoi(arr[1])
	}

	myExp = regexp.MustCompile(`[\S\s]*title=([^\n]*?)\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		title = strings.TrimSpace(arr[1])
	}

	if caseId <= 0 {
		pass = false
	}

	return
}

func ReadExpectIndependentArr(content string) [][]string {
	//正常显示6
	//E2.16
	//>>
	//  E2.2 - 16
	//  E2.2 - 26
	//>>
	//>>
	//  E3 - 16
	//  E3 - 26
	//>>

	lines := strings.Split(content, "\n")

	ret := make([][]string, 0)
	var cpArr []string

	currModel := ""
	idx := 0
	for idx < len(lines) {
		line := strings.TrimSpace(lines[idx])

		if line == ">>" { // more than one line
			currModel = "multi"
			cpArr = make([]string, 0)
		} else if currModel == "multi" { // in >> and >> in multi line mode
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.Index(lines[idx+1], ">>") > -1 { // end multi line
				temp := make([]string, 0)
				temp = append(temp, strings.Join(cpArr, "\r\n"))

				ret = append(ret, temp)
				cpArr = make([]string, 0)
				currModel = ""

				idx += 1
			}
		} else {
			currModel = "single"

			line = strings.TrimSpace(line)

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		}

		idx += 1
	}

	return ret
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
				temp = append(temp, cpArr...)

				ret = append(ret, temp)
				cpArr = make([]string, 0)

				idx = idx + 1
				model = ""
			}
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
	pass, _ := regexp.MatchString(`cid\b\s*=`, content)

	return pass
}

func ReadCaseInfo(content, lang string) (info, checkpoints string) {
	regStr := fmt.Sprintf(`(?smU)%s((?U:.*pid.*))\n(.*)%s`,
		commConsts.LangCommentsRegxMap[lang][0], commConsts.LangCommentsRegxMap[lang][1])

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
	dir := fileUtils.AddFilePathSepIfNeeded(filepath.Dir(file))
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

func GetScriptByIdsInDir(dirPth string, idMap *map[int]string) error {
	dirPth = fileUtils.AbsolutePath(dirPth)

	sep := consts.FilePthSep

	if commonUtils.IgnoreZtfFile(dirPth) {
		return nil
	}

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			GetScriptByIdsInDir(dirPth+name+sep, idMap)
		} else {
			regx := langHelper.GetSupportLanguageExtRegx()
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", name)

			if !pass {
				continue
			}

			path := dirPth + name
			pass, id, _, _, _ := GetCaseInfo(path)
			if pass {
				(*idMap)[id] = path
			}
		}
	}

	return nil
}

func GetCaseIdsInSuiteFile(name string, ids *[]int) {
	content := fileUtils.ReadFile(name)

	for _, line := range strings.Split(content, "\n") {
		idStr := strings.TrimSpace(line)
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err == nil {
			*ids = append(*ids, id)
		}
	}
}

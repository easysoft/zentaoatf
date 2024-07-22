package scriptHelper

import (
	"fmt"
	"os"
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
	"github.com/ergoapi/util/file"
)

func GetStepAndExpectMap(file string) (steps []commDomain.ZentaoCaseStep) {
	if !fileUtils.FileExist(file) {
		return
	}

	lang := langHelper.GetLangByFile(file)
	content := fileUtils.ReadFile(file)

	info, checkpoints := ReadCaseInfoInOldFormat(content, lang)
	if info != "" {
		steps = GetStepAndExpectMapInOldFormat(checkpoints, file)
		return
	}

	_, _, steps = ReadTitleAndStepsInNewFormat(content, lang)

	return
}

func ReadTitleAndStepsInNewFormat(content, lang string) (caseId int, title string, steps []commDomain.ZentaoCaseStep) {
	//测试用例标题 #1
	//- 步骤1 @期待结果1
	//- 步骤2
	//- 子步骤2.1 @{
	//	期待结果2.1.1
	//	期待结果2.1.2
	//}
	//- 子步骤2.2 @期待结果2.2
	//- 步骤3 @期待结果3

	comments := strings.TrimSpace(getScriptComments(content, lang))
	index := 0
	groupName := 1
	itemName := 1
	titleLineStart := false
	lines := strings.Split(comments, "\n")
	for index < len(lines) {
		line := lines[index]

		if !titleLineStart {
			caseId, title = findTitle(line)
			if title != "" {
				titleLineStart = true
				index += 1
				continue
			}
		}

		isStepLine, descAndExpect, isChild := isStepLine(line)
		if !isStepLine {
			index += 1
			continue
		}

		isMultiExpect, desc2 := isMultiLineExpectStart(line)
		step := commDomain.ZentaoCaseStep{}

		if isMultiExpect { // more than one line
			index += 1
			step.Desc = desc2
			step.Expect = getMultiExpect(lines, &index)
		} else {
			step.Desc, step.Expect = getSingleExpect(descAndExpect)
		}

		step.Type = commConsts.Group
		if isChild {
			step.Type = commConsts.Item
			step.Name = fmt.Sprintf("%d.%d", groupName-1, itemName)
			itemName += 1
		} else {
			step.Name = fmt.Sprintf("%d", groupName)
			groupName += 1
			itemName = 1
		}

		steps = append(steps, step)

		index += 1
	}

	return
}

func getSingleExpect(descAndExpect string) (desc, expect string) {
	arr := strings.Split(descAndExpect, "@")

	desc = strings.TrimSpace(arr[0])
	if len(arr) > 1 {
		expect = strings.TrimSpace(arr[1])
		if expect == "" {
			expect = "pass"
		}
	}

	return
}

func findTitle(line string) (id int, title string) {
	reg := `(.*)#(\d*)`
	arr := regexp.MustCompile(reg).FindStringSubmatch(line)
	if len(arr) > 2 {
		var err error
		id, err = strconv.Atoi(arr[2])
		if err == nil {
			title = strings.TrimSpace(arr[1])
		}
	}

	return
}

func getMultiExpect(lines []string, index *int) (ret string) {
	var arr []string

	for *index < len(lines) {
		line := strings.TrimSpace(lines[*index])
		if isMultiLineExpectEnd(line) {
			break
		}

		arr = append(arr, line)

		*index += 1
	}

	ret = strings.Join(arr, "\r\n")
	return
}

func isStepLine(line string) (is bool, ret string, isChild bool) {
	reg := `^(\s*)-\s*(.+)$`
	arr := regexp.MustCompile(reg).FindStringSubmatch(line)
	if len(arr) > 2 {
		is = true
		ret = arr[2]

		if len(arr[1]) > 1 {
			isChild = true
		}
	}

	return
}
func isMultiLineExpectStart(line string) (is bool, step string) {
	reg := `^\s*-\s*(.+)@\s*\{\s*$`
	arr := regexp.MustCompile(reg).FindStringSubmatch(line)
	if len(arr) > 1 {
		is = true
		step = arr[1]
	}

	return
}
func isMultiLineExpectEnd(line string) (is bool) {
	is = strings.TrimSpace(line) == "}"
	return
}

func ReadCaseInfoInOldFormat(content, lang string) (info, checkpoints string) {
	regStr := fmt.Sprintf(`(?smU)%s((?U:.*pid\s*=.*))\n(.*)%s`,
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

func GetStepAndExpectMapInOldFormat(checkpoints, file string) (steps []commDomain.ZentaoCaseStep) {
	lines := strings.Split(checkpoints, "\n")

	groupArr := getStepNestedArrInOldFormat(lines)
	_, steps = getSortedTextFromNestedSteps(groupArr)

	isIndependent, expectIndependentContent := GetDependentExpect(file)
	if isIndependent {
		GetExpectMapFromIndependentFile(&steps, expectIndependentContent, false)
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

func getStepNestedArrInOldFormat(lines []string) (ret []commDomain.ZtfStep) {
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

			if strings.TrimSpace(parent.Expect) == "" && strings.Contains(line, ">>") {
				parent.Expect = commConsts.ExpectResultPass
			}
			ret = append(ret, parent)

			continue
		}

		// 有缩进
		child := commDomain.ZtfStep{}
		child, increase = parserNextLines(line, lines[index+1:])
		index += increase

		if parent.Desc == "" {
			continue
		}

		if strings.TrimSpace(child.Expect) == "" && strings.Contains(line, ">>") {
			child.Expect = commConsts.ExpectResultPass
		}

		ret[len(ret)-1].Children = append(ret[len(ret)-1].Children, child)
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

	if !strings.Contains(str, ">>") || expect != "" { // no >> or single line expect
		ret = commDomain.ZtfStep{Desc: desc, Expect: expect}
		return
	}

	// will test if it has multi-line expect
	for index, line := range nextLines {
		if strings.TrimSpace(line) == ">>" {
			increase = index
			break
		}

		if strings.Contains(line, ">>") {
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

	ret = commDomain.ZtfStep{Desc: desc, Expect: expect}
	return
}

func ReplaceCaseDesc(desc, file string) {
	content := fileUtils.ReadFile(file)
	lang := langHelper.GetLangByFile(file)

	regStr := fmt.Sprintf(`(?smU)%s((?U:.*cid.*))\n(.*)%s`,
		commConsts.LangCommentsRegxMap[lang][0], commConsts.LangCommentsRegxMap[lang][1])
	re, _ := regexp.Compile(regStr)

	newDesc := fmt.Sprintf("%s\n"+desc+"\n%s",
		commConsts.LangCommentsTagMap[lang][0],
		commConsts.LangCommentsTagMap[lang][1])
	newDesc = strings.Replace(newDesc, "$", "￥￥￥", -1)

	out := re.ReplaceAllString(content, newDesc)

	out = strings.Replace(out, "￥￥￥", "$", -1)
	fileUtils.WriteFile(file, out)
}

func getSortedTextFromNestedSteps(groups []commDomain.ZtfStep) (ret string, steps []commDomain.ZentaoCaseStep) {
	arr := make([]string, 0)

	for _, group := range groups {
		step := commDomain.ZentaoCaseStep{}

		stepType := commConsts.Group
		step.Type = stepType

		stepTxt := strings.TrimSpace(group.Desc)
		step.Desc = stepTxt

		expectTxt := strings.TrimSpace(group.Expect)
		expectTxt = strings.TrimSuffix(expectTxt, "]]")
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

func ScriptToExpectName(file string) string {
	fileSuffix := path.Ext(file)
	expectName := strings.TrimSuffix(file, fileSuffix) + ".exp"

	return expectName
}

func getScriptComments(content, lang string) (ret string) {
	if lang == "" {
		return
	}
	reg := fmt.Sprintf(`(?smU)%s((?U:.*))%s`, commConsts.LangCommentsRegxMap[lang][0], commConsts.LangCommentsRegxMap[lang][1])
	arr := regexp.MustCompile(reg).FindStringSubmatch(content)
	if len(arr) < 2 { // wrong format
		return
	}

	ret = strings.TrimSpace(arr[1])

	return
}

func GetCaseInfo(file string) (pass bool, caseId, productId int, title string, timeout int64) {
	content := fileUtils.ReadFile(file)
	lang := langHelper.GetLangByFile(file)

	comments := strings.TrimSpace(getScriptComments(content, lang))
	index := 0
	lines := strings.Split(comments, "\n")
	for index < len(lines) {
		line := strings.TrimSpace(lines[index])
		caseId, title = findTitle(line)
		if title != "" {
			break
		}

		index += 1
	}

	pass = title != ""
	if pass {
		return
	}

	// deal with old format, will be removed
	isOldFormat := strings.Index(content, "[esac]") > -1
	pass = CheckFileContentIsScript(content)
	if !pass {
		return false, caseId, productId, title, timeout
	}

	caseInfo := ""
	regStr := ""
	if isOldFormat {
		regStr = `(?s)\[case\](.*)\[esac\]`
	} else {
		regStr = fmt.Sprintf(`(?sm)%s((?U:.*cid.*))\n(.*)%s`,
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

	myExp = regexp.MustCompile(`[\S\s]*cid=\s*([^\n]*?)\s*\n`)
	arr = myExp.FindStringSubmatch(caseInfo)
	if len(arr) > 1 {
		productId, _ = strconv.Atoi(arr[1])
	}

	myExp = regexp.MustCompile(`[\S\s]*title\s*=\s*([^\n]*?)\n`)
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
	//{
	//  E2.2 - 16
	//  E2.2 - 26
	//}
	//{
	//  E3 - 16
	//  E3 - 26
	//}

	lines := strings.Split(content, "\n")

	ret := make([][]string, 0)
	var cpArr []string

	currModel := ""
	idx := 0
	for idx < len(lines) {
		line := strings.TrimSpace(lines[idx])

		if line == "{" { // more than one line
			currModel = "multi"
			cpArr = make([]string, 0)
		} else if currModel == "multi" { // in { and } in multi line mode
			cpArr = append(cpArr, line)

			if idx == len(lines)-1 || strings.TrimSpace(lines[idx+1]) == "}" { // end multi line
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

		if line == "@{" { // more than one line
			model = "multi"
			cpArr = make([]string, 0)
		} else if model == "multi" { // between line @{ and } in multi line mode
			if line != "" {
				cpArr = append(cpArr, line)
			}

			//if idx == len(lines)-1 || strings.Index(lines[idx+1], "}") > -1 {
			if idx == len(lines)-1 || strings.TrimSpace(lines[idx+1]) == "}" {
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
			if strings.Index(line, "@") != 0 { // ignore the line not started with @
				continue
			}

			line = line[1:]

			cpArr = append(cpArr, line)
			ret = append(ret, cpArr)
			cpArr = make([]string, 0)
		}
	}

	return
}

func ReadLogArrOld(content string) (isSkip bool, ret [][]string) {
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
			if line != "" {
				cpArr = append(cpArr, line)
			}

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
	os.Chmod(path, 0777)
	content := fileUtils.ReadFile(path)

	pass := CheckFileContentIsScript(content)
	return pass
}

func CheckFileContentIsScript(content string) bool {
	pass, _ := regexp.MatchString(`cid\b\s*=`, content)

	if !pass {
		pass, _ = regexp.MatchString(`(?m:^(.+ +)#\d*$)`, content)
	}

	return pass
}

func GetScriptByIdsInDir(dirPth string, idMap *map[int]string) error {
	if fileUtils.IsFile(dirPth) {
		pass, id, _, _, _ := GetCaseInfo(dirPth)
		if pass {
			(*idMap)[id] = dirPth
		}
		return nil
	}
	dirPth = fileUtils.AbsolutePath(dirPth)

	sep := consts.FilePthSep

	if commonUtils.IgnoreZtfFile(dirPth) {
		return nil
	}

	dir, err := file.ReadDir(dirPth)
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
			} else {
				pass, id, _, _, _ = GetCaseInfo(path)
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

//func getGroupBlockArr(lines []string) [][]string {
//	groupBlockArr := make([][]string, 0)
//
//	idx := 0
//	for true {
//		if idx >= len(lines) {
//			break
//		}
//
//		var groupContent []string
//		line := strings.TrimSpace(lines[idx])
//		if isGroup(line) { // must match a group
//			groupContent = make([]string, 0)
//			groupContent = append(groupContent, line)
//
//			idx++
//
//			for true {
//				if idx >= len(lines) {
//					groupBlockArr = append(groupBlockArr, groupContent)
//					break
//				}
//
//				line = strings.TrimSpace(lines[idx])
//				if isGroup(line) {
//					groupBlockArr = append(groupBlockArr, groupContent)
//
//					break
//				} else if line != "" && !isGroup(line) {
//					groupContent = append(groupContent, line)
//				}
//
//				idx++
//			}
//		} else {
//			idx++
//		}
//	}
//
//	return groupBlockArr
//}

//func loadMultiLineSteps(arr []string) []commDomain.ZtfStep {
//	childs := make([]commDomain.ZtfStep, 0)
//
//	child := commDomain.ZtfStep{}
//	idx := 0
//	for true {
//		if idx >= len(arr) {
//			if child.Desc != "" {
//				childs = append(childs, child)
//			}
//
//			break
//		}
//
//		line := arr[idx]
//		line = strings.TrimSpace(line)
//
//		if isStepsIdent(line) {
//			if idx > 0 {
//				childs = append(childs, child)
//			}
//
//			child = commDomain.ZtfStep{}
//			idx++
//
//			stp := ""
//			for true { // retrieve next lines
//				if idx >= len(arr) || hasBrackets(arr[idx]) {
//					child.Desc = stp
//					break
//				}
//
//				stp += arr[idx] + "\n"
//				idx++
//			}
//		}
//
//		if isExpectsIdent(line) {
//			idx++
//
//			exp := ""
//			for true { // retrieve next lines
//				if idx >= len(arr) || hasBrackets(arr[idx]) {
//					child.Expect = exp
//					break
//				}
//
//				temp := strings.TrimSpace(arr[idx])
//				if temp == ">>" {
//					temp = ""
//				}
//				exp += temp + "\n"
//				idx++
//			}
//		}
//
//	}
//
//	return childs
//}
//
//func loadSingleLineSteps(arr []string) []commDomain.ZtfStep {
//	children := make([]commDomain.ZtfStep, 0)
//
//	for _, line := range arr {
//		line = strings.TrimSpace(line)
//
//		sections := strings.Split(line, ">>")
//		expect := ""
//		if len(sections) > 1 { // has expect
//			expect = strings.TrimSpace(sections[1])
//		}
//
//		child := commDomain.ZtfStep{Desc: sections[0], Expect: expect}
//
//		children = append(children, child)
//	}
//
//	return children
//}
//
//func isGroupIdent(str string) bool {
//	pass, _ := regexp.MatchString(`(?i)\[\s*group\s*\]`, str)
//	return pass
//}

//func isGroup(str string) bool {
//	ret := strings.Index(str, ">>") < 0 && hasBrackets(str) && !isStepsIdent(str) && !isExpectsIdent(str)
//
//	return ret
//}

//func isStepsIdent(str string) bool {
//	pass, _ := regexp.MatchString(`(?i)\[.*steps\.*\]`, str)
//	return pass
//}
//
//func isExpectsIdent(str string) bool {
//	pass, _ := regexp.MatchString(`(?i)\[.*expects\.*\]`, str)
//	return pass
//}
//
//func hasBrackets(str string) bool {
//	pass, _ := regexp.MatchString(`(?i)()\[.*\]`, str)
//	return pass
//}

//func getGroupName(str string) string {
//	reg := `\[\d\.\s]*(.*)\]`
//	repl := "${1}"
//
//	regx, _ := regexp.Compile(reg)
//	str = regx.ReplaceAllString(str, repl)
//
//	return str
//}
//
//func printMultiStepOrExpect(str string) string {
//	str = strings.TrimSpace(str)
//
//	ret := make([]string, 0)
//
//	for _, line := range strings.Split(str, "\n") {
//		line = strings.TrimSpace(line)
//
//		ret = append(ret, fmt.Sprintf("%s%s", strings.Repeat(" ", 4), line))
//	}
//
//	return strings.Join(ret, "\r\n")
//}

//func replaceNumb(str string, groupNumb int, childNumb int, withBrackets bool) string {
//	numb := getNumbStr(groupNumb, childNumb)
//
//	reg := `[\d\.\s]*(.*)`
//	repl := numb + " ${1}"
//	if withBrackets {
//		reg = `\[` + reg + `\]`
//		repl = `[` + repl + `]`
//	}
//
//	regx, _ := regexp.Compile(reg)
//	str = regx.ReplaceAllString(str, repl)
//
//	return str
//}
//func getNumbStr(groupNumb int, childNumb int) string {
//	numb := strconv.Itoa(groupNumb) + "."
//	if childNumb != -1 {
//		numb += strconv.Itoa(childNumb) + "."
//	}
//
//	return numb
//}

//func IsMultiLine(step commDomain.ZtfStep) bool {
//	if strings.Index(step.Desc, "\n") > -1 || strings.Index(step.Expect, "\n") > -1 {
//		return true
//	}
//
//	return false
//}
//func RunDateFolder() string {
//	runName := dateUtils.DateTimeStrFmt(time.Now(), "2006-01-02T150405") + string(os.PathSeparator)
//
//	return runName
//}

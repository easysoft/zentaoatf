package action

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
	"strconv"
	"strings"
)

func Gen(remoteUrl string, langType string, independentExpectFile bool) {
	buf := utils.ReadFileBuf(remoteUrl)

	var resp model.Response
	json.Unmarshal(buf, &resp)

	if resp.Code != 1 {
		fmt.Println(string(buf))
		return
	}

	for _, testCase := range resp.Cases {
		DealwithTestCase(testCase, langType, independentExpectFile)
	}
}

func DealwithTestCase(tc model.TestCase, langType string, independentExpect bool) {
	langMap := GetLangMap()
	StepWidth := 20

	caseId := tc.Id
	caseTitle := tc.Title
	scriptFile := fmt.Sprintf("xdoc/scripts/tc-%s.%s", strconv.Itoa(caseId), langMap[langType]["extName"])

	steps := make([]string, 0)
	expects := make([]string, 0)
	srcCode := make([]string, 0)

	steps = append(steps, "@开头的为含验证点的步骤")
	temp := fmt.Sprintf("\n%sCODE: 此处编写操作步骤代码\n", langMap[langType]["commentsTag"])
	srcCode = append(srcCode, temp)
	readme := utils.ReadFile("xdoc/template/readme.tpl") + "\n"

	stepDisplayMaxWidth := 0
	DealwithTestStepWidth(tc.Steps, &stepDisplayMaxWidth, StepWidth)

	level := 1
	checkPointIndex := 0
	for _, ts := range tc.Steps {
		DealwithTestStep(ts, langType, level, StepWidth, &checkPointIndex, &steps, &expects, &srcCode)
	}

	var expectsTxt string
	if independentExpect {
		expectFile := utils.ScriptToExpectName(scriptFile)

		expectsTxt = "@file"
		utils.WriteFile(expectFile, strings.Join(expects, "\n"))
	} else {
		expectsTxt = strings.Join(expects, "\n")
	}

	template := utils.ReadFile("xdoc/template/" + langType + ".tpl")
	content := fmt.Sprintf(template,
		caseId, caseTitle,
		strings.Join(steps, "\n"), expectsTxt,
		readme,
		strings.Join(srcCode, "\n"))

	fmt.Println(content)

	utils.WriteFile(scriptFile, content)
}

func DealwithTestStepWidth(steps []model.TestStep, stepSDisplayMaxWidth *int, stepWidth int) {
	for _, ts := range steps {
		length := len(strconv.Itoa(ts.Id))
		if length > *stepSDisplayMaxWidth {
			*stepSDisplayMaxWidth = length
		}
	}
	*stepSDisplayMaxWidth += stepWidth // prefix space and @step
}

func DealwithTestStep(ts model.TestStep, langType string,
	level int, stepWidth int, checkPointIndex *int,
	steps *[]string, expects *[]string, srcCode *[]string) {
	langMap := GetLangMap()

	isGroup := ts.IsGroup
	isCheckPoint := ts.IsCheckPoint

	stepId := ts.Id
	stepTitle := ts.Title
	stepExpect := ts.Expect

	// 处理steps
	var stepType string
	if isGroup {
		stepType = "group"
	} else {
		stepType = "step"
	}

	stepIdent := stepType + strconv.Itoa(stepId)
	if isCheckPoint {
		stepIdent = "@" + stepIdent
		*checkPointIndex++
	}

	preFixSpace := level * 3
	postFixSpace := stepWidth - preFixSpace - len(stepIdent)

	stepLine := fmt.Sprintf("%*s", preFixSpace, " ") + stepIdent
	stepLine += fmt.Sprintf("%*s", postFixSpace, " ")
	stepLine += stepTitle

	*steps = append(*steps, stepLine)

	// 处理expects
	if isCheckPoint {
		expectsLine := ""

		expectsLine = "# \n"
		expectsLine += "CODE: " + stepIdent + "期望结果, 可以有多行\n"

		*expects = append(*expects, expectsLine)
	}

	// 处理srcCode
	if isCheckPoint {
		codeLine := langMap[langType]["printGrammar"]

		codeLine += fmt.Sprintf("  %s %s: %s\n", langMap[langType]["commentsTag"], stepIdent, stepExpect)

		codeLine += langMap[langType]["commentsTag"] + "CODE: 输出验证点实际结果\n"

		*srcCode = append(*srcCode, codeLine)
	}

	if isGroup {
		for _, tsChild := range ts.Steps {
			DealwithTestStep(tsChild, langType, level+1, stepWidth, checkPointIndex, steps, expects, srcCode)
		}
	}
}

func GetLangMap() map[string]map[string]string {
	langMap := make(map[string]map[string]string)

	langMap["go"] = map[string]string{
		"extName":      "go",
		"commentsTag":  "//",
		"printGrammar": "println(\"#\")",
	}

	langMap["lua"] = map[string]string{
		"extName":      "lua",
		"commentsTag":  "--",
		"printGrammar": "print('#')",
	}

	langMap["perl"] = map[string]string{
		"extName":      "pl",
		"commentsTag":  "#",
		"printGrammar": "print \"#\\n\";",
	}

	langMap["php"] = map[string]string{
		"extName":      "php",
		"commentsTag":  "//",
		"printGrammar": "echo \"#\n\";",
	}

	langMap["python"] = map[string]string{
		"extName":      "py",
		"commentsTag":  "#",
		"printGrammar": "print(\"#\")",
	}

	langMap["ruby"] = map[string]string{
		"extName":      "rb",
		"commentsTag":  "#",
		"printGrammar": "print(\"#\\n\")",
	}

	langMap["shell"] = map[string]string{
		"extName":      "sh",
		"commentsTag":  "#",
		"printGrammar": "echo \"#\"",
	}

	langMap["tcl"] = map[string]string{
		"extName":      "tl",
		"commentsTag":  "#",
		"printGrammar": "set hello \"#\"; \n puts [set hello];",
	}

	return langMap
}

package action

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zentaoatf/src/misc"
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
	StepWidth := 20

	caseId := tc.Id
	caseTitle := tc.Title
	scriptFile := "xdoc/scripts/tc-" + strconv.Itoa(caseId) + "." + langType

	steps := make([]string, 0)
	expects := make([]string, 0)
	srcCode := make([]string, 0)

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

	template := utils.ReadFile("xdoc/script-template.txt")
	content := string(fmt.Sprintf(string(template), langType, caseId, caseTitle,
		strings.Join(steps, "\n"), expectsTxt, strings.Join(srcCode, "\n")))

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

func DealwithTestStep(ts model.TestStep, langType string, level int, stepWidth int, checkPointIndex *int,
	steps *[]string, expects *[]string, srcCode *[]string) {

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
		expectsLine += "/* " + stepIdent + "期望结果, 可以有多行 */\n"

		*expects = append(*expects, expectsLine)
	}

	// 处理srcCode
	if isCheckPoint {
		codeLine := ""

		if langType == misc.PHP.String() {
			codeLine += `echo "#\n";`
		} else if langType == misc.GO.String() {
			codeLine += `println("#")\n`
		}

		codeLine += "  // " + stepIdent + ": " + stepExpect + "\n"
		codeLine += "/* 输出验证点实际结果 */\n"

		*srcCode = append(*srcCode, codeLine)
	}

	if isGroup {
		for _, tsChild := range ts.Steps {
			DealwithTestStep(tsChild, langType, level+1, stepWidth, checkPointIndex, steps, expects, srcCode)
		}
	}
}

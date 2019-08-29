package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/lang"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"os"
	"regexp"
	"strings"
)

func Generate(testcases []model.TestCase, langType string, independentFile bool) (int, error) {
	caseIds := make([]string, 0)
	for _, cs := range testcases {
		GenerateTestCaseScript(cs, langType, independentFile, &caseIds)
	}

	GenSuite(caseIds)

	return len(testcases), nil
}

func GenerateTestCaseScript(cs model.TestCase, langType string, independentFile bool, caseIds *[]string) {
	caseId := cs.Id
	productId := cs.Product
	moduleId := cs.Module
	caseTitle := cs.Title

	modulePath := ""
	if vari.ZentaoCaseFileds.Modules[moduleId] != "" {
		modulePath = vari.ZentaoCaseFileds.Modules[moduleId] + string(os.PathSeparator)
		modulePath = modulePath[1:]
	}

	scriptFile := fmt.Sprintf(constant.ScriptDir+"%stc-%s.%s", modulePath, caseId, langUtils.LangMap[langType]["extName"])

	fileUtils.MkDirIfNeeded(constant.ScriptDir)
	*caseIds = append(*caseIds, caseId)

	steps := make([]string, 0)
	expects := make([]string, 0)
	srcCode := make([]string, 0)

	steps = append(steps, i118Utils.I118Prt.Sprintf("is_checkpoint"))

	StepWidth := 20
	stepDisplayMaxWidth := 0
	computerTestStepWidth(cs.StepArr, &stepDisplayMaxWidth, StepWidth)

	for _, ts := range cs.StepArr {
		GenerateTestStepScript(ts, langType, StepWidth, &steps, &expects, &srcCode)
	}

	if fileUtils.FileExist(scriptFile) { // update title and steps
		content := fileUtils.ReadFile(scriptFile)

		// replace title
		re, _ := regexp.Compile(`title:\s*([^\n]*?)\s*\n`)
		content = re.ReplaceAllString(content, fmt.Sprintf("title:          %s\n", caseTitle))

		// replace steps
		re, _ = regexp.Compile(`steps:[^\n]*\n*([\S\s]*)\n+expects:`)
		content = re.ReplaceAllString(content,
			fmt.Sprintf("steps:          %s\n\nexpects:", strings.Join(steps, "\n")))

		fileUtils.WriteFile(scriptFile, content)
		return
	}

	temp := fmt.Sprintf("\n%sCODE: %s", langUtils.LangMap[langType]["commentsTag"],
		i118Utils.I118Prt.Sprintf("your_codes_here"))
	srcCode = append(srcCode, temp)

	var expectsTxt string
	if !independentFile {
		expectsTxt = strings.Join(expects, "\n")
	} else {
		expectFile := zentaoUtils.ScriptToExpectName(scriptFile)

		expectsTxt = "@file\n"
		fileUtils.WriteFile(expectFile, strings.Join(expects, "\n"))
	}

	path := fmt.Sprintf("res%stemplate%s", string(os.PathSeparator), string(os.PathSeparator))
	template := zentaoUtils.ReadResData(path + langType + ".tpl")

	content := fmt.Sprintf(template,
		caseId, productId, caseTitle,
		strings.Join(steps, "\n"), expectsTxt,
		strings.Join(srcCode, "\n"))

	fileUtils.WriteFile(scriptFile, content)
}

func GenerateTestStepScript(ts model.TestStep, langType string, stepWidth int,
	steps *[]string, expects *[]string, srcCode *[]string) {
	LangMap := langUtils.GetSupportedScriptLang()

	isGroup := ts.Type == "group"
	isCheckPoint := ts.Expect != ""

	stepId := ts.Id
	stepTitle := ts.Desc
	stepExpect := ts.Expect
	stepParent := ts.Parent

	// 处理steps
	preFixSpace := 3
	if stepParent != "" && stepParent != "0" {
		preFixSpace = 6
	}

	var stepType string
	if isGroup {
		stepType = "group"
	} else {
		stepType = "step"
	}

	stepIdent := stepType + stepId
	if isCheckPoint {
		stepIdent = "@" + stepIdent
	}

	postFixSpace := stepWidth - preFixSpace - len(stepIdent)

	stepLine := fmt.Sprintf("%*s", preFixSpace, " ") + stepIdent
	stepLine += fmt.Sprintf("%*s", postFixSpace, " ")
	stepLine += stepTitle

	*steps = append(*steps, stepLine)

	// 处理expects
	if isCheckPoint {
		expectsLine := ""

		expectsLine = "# " + stepIdent + " \n"
		expectsLine += "CODE: " + i118Utils.I118Prt.Sprintf("expect_result_here") + " \n"

		*expects = append(*expects, expectsLine)
	}

	// 处理srcCode
	if isCheckPoint {
		codeLine := LangMap[langType]["printGrammar"]

		codeLine += fmt.Sprintf("  %s %s: %s\n", LangMap[langType]["commentsTag"], stepIdent, stepExpect)

		codeLine += LangMap[langType]["commentsTag"] + "CODE: " + i118Utils.I118Prt.Sprintf("actual_result_here")

		*srcCode = append(*srcCode, codeLine)
	}
}

func GenSuite(cases []string) {
	str := strings.Join(cases, "\n")

	fileUtils.WriteFile(constant.ScriptDir+"all."+constant.ExtNameSuite, str)
}

func computerTestStepWidth(steps []model.TestStep, stepSDisplayMaxWidth *int, stepWidth int) {
	for _, ts := range steps {
		length := len(ts.Id)
		if length > *stepSDisplayMaxWidth {
			*stepSDisplayMaxWidth = length
		}
	}
	*stepSDisplayMaxWidth += stepWidth // prefix space and @step
}

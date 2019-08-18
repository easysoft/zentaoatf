package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	dateUtils "github.com/easysoft/zentaoatf/src/utils/date"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"os"
	"strings"
	"time"
)

func Generate(testcases []model.TestCase, langType string, singleFile bool,
	account string, password string) (int, error) {

	casePaths := make([]string, 0)
	for _, cs := range testcases {
		GenerateTestCaseScrrpt(cs, langType, singleFile, &casePaths)
	}

	GenSuite(casePaths)

	return len(testcases), nil
}

func GenerateTestCaseScrrpt(cs model.TestCase, langType string, singleFile bool, casePaths *[]string) {
	LangMap := GetSupportedScriptLang()
	langs := ""
	if LangMap[langType] == nil {
		i := 0
		for lang, _ := range LangMap {
			if i > 0 {
				langs += ", "
			}
			langs += lang
			i++
		}
		fmt.Printf("only support languages %s \n", langs)
		os.Exit(1)
	}

	caseId := cs.Id
	caseInTaskId := cs.IdInTask
	taskId := cs.TaskId
	if caseInTaskId == "" {
		caseInTaskId = "0"
	}
	caseTitle := cs.Title

	scriptFile := fmt.Sprintf(constant.ScriptDir+"tc-%s.%s", caseId, LangMap[langType]["extName"])
	if fileUtils.FileExist(scriptFile) {
		scriptFile = fmt.Sprintf(constant.ScriptDir+"tc-%s.%s",
			caseId+"-"+dateUtils.DateTimeStrLong(time.Now()), LangMap[langType]["extName"])
	}

	fileUtils.MkDirIfNeeded(vari.Prefer.WorkDir + constant.ScriptDir)
	*casePaths = append(*casePaths, scriptFile)
	scriptFullPath := vari.Prefer.WorkDir + scriptFile

	steps := make([]string, 0)
	expects := make([]string, 0)
	srcCode := make([]string, 0)

	steps = append(steps, "@开头的为含验证点的步骤")

	temp := fmt.Sprintf("\n%sCODE: 此处编写操作步骤代码\n", LangMap[langType]["commentsTag"])
	srcCode = append(srcCode, temp)

	readme := zentaoUtils.ReadResData("res/template/readme.tpl") + "\n"

	StepWidth := 20
	stepDisplayMaxWidth := 0
	computerTestStepWidth(cs.StepArr, &stepDisplayMaxWidth, StepWidth)

	for _, ts := range cs.StepArr {
		GenerateTestStepScript(ts, langType, StepWidth, &steps, &expects, &srcCode)
	}

	var expectsTxt string
	if singleFile {
		expectsTxt = strings.Join(expects, "\n")
	} else {
		expectFile := zentaoUtils.ScriptToExpectName(scriptFullPath)

		expectsTxt = "@file\n"
		fileUtils.WriteFile(expectFile, strings.Join(expects, "\n"))
	}

	path := fmt.Sprintf("res%stemplate%s", string(os.PathSeparator), string(os.PathSeparator))
	template := zentaoUtils.ReadResData(path + langType + ".tpl")

	content := fmt.Sprintf(template,
		caseId, caseInTaskId, taskId, caseTitle,
		strings.Join(steps, "\n"), expectsTxt,
		readme,
		strings.Join(srcCode, "\n"))

	//fmt.Println(content)

	fileUtils.WriteFile(scriptFullPath, content)
}

func GenerateTestStepScript(ts model.TestStep, langType string, stepWidth int,
	steps *[]string, expects *[]string, srcCode *[]string) {
	LangMap := GetSupportedScriptLang()

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
		expectsLine += "CODE: " + "期望结果, 可以有多行\n"

		*expects = append(*expects, expectsLine)
	}

	// 处理srcCode
	if isCheckPoint {
		codeLine := LangMap[langType]["printGrammar"]

		codeLine += fmt.Sprintf("  %s %s: %s\n", LangMap[langType]["commentsTag"], stepIdent, stepExpect)

		codeLine += LangMap[langType]["commentsTag"] + "CODE: 输出验证点实际结果\n"

		*srcCode = append(*srcCode, codeLine)
	}
}

func GenSuite(cases []string) {
	str := strings.Join(cases, "\n")

	fileUtils.WriteFile(vari.Prefer.WorkDir+constant.ScriptDir+"all."+constant.ExtNameSuite, str)
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

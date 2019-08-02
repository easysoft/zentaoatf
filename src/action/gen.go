package action

import (
	"errors"
	"fmt"
	"github.com/easysoft/zentaoatf/src/biz"
	httpClient "github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenFromCmd(url string, entityType string, entityVal string, langType string, singleFile bool) {
	params := make(map[string]string)

	params["entityType"] = entityType
	params["entityVal"] = entityVal

	json, err := httpClient.Get(url, params)

	if err == nil {
		Generate(json, url, entityType, entityVal, langType, singleFile)
	}
}

func Generate(json model.Response,
	url string, entityType string, entityVal string, langType string, singleFile bool) (int, error) {
	if json.Code != 1 {
		return 0, errors.New("response code = %s")
	}

	casePaths := make([]string, 0)
	for _, testCase := range json.Cases {
		DealwithTestCase(testCase, langType, singleFile, &casePaths)
	}
	biz.GenSuite(casePaths)

	utils.SaveConfig("", url, entityType, entityVal, langType, singleFile, json.Name)

	return len(json.Cases), nil
}

func DealwithTestCase(tc model.TestCase, langType string, singleFile bool, casePaths *[]string) {
	LangMap := script.GetLangMap()
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

	StepWidth := 20

	caseId := tc.Id
	caseTitle := tc.Title

	scriptFile := fmt.Sprintf(utils.GenDir+"tc-%s.%s", strconv.Itoa(caseId), LangMap[langType]["extName"])

	utils.MkDirIfNeeded(utils.Prefer.WorkDir + utils.GenDir)
	if utils.FileExist(utils.Prefer.WorkDir + scriptFile) {
		scriptFile = fmt.Sprintf(utils.GenDir+"tc-%s.%s",
			strconv.Itoa(caseId)+"-"+utils.DateTimeStrLong(time.Now()), LangMap[langType]["extName"])
	}
	*casePaths = append(*casePaths, scriptFile)
	scriptFullPath := utils.Prefer.WorkDir + scriptFile

	steps := make([]string, 0)
	expects := make([]string, 0)
	srcCode := make([]string, 0)

	steps = append(steps, "@开头的为含验证点的步骤")

	temp := fmt.Sprintf("\n%sCODE: 此处编写操作步骤代码\n", LangMap[langType]["commentsTag"])
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
	if singleFile {
		expectsTxt = strings.Join(expects, "\n")
	} else {
		expectFile := utils.ScriptToExpectName(scriptFullPath)

		expectsTxt = "@file\n"
		utils.WriteFile(expectFile, strings.Join(expects, "\n"))
	}

	template := utils.ReadFile("xdoc/template/" + langType + ".tpl")
	content := fmt.Sprintf(template,
		caseId, caseTitle,
		strings.Join(steps, "\n"), expectsTxt,
		readme,
		strings.Join(srcCode, "\n"))

	//fmt.Println(content)

	utils.WriteFile(scriptFullPath, content)
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
	LangMap := script.GetLangMap()

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
		codeLine := LangMap[langType]["printGrammar"]

		codeLine += fmt.Sprintf("  %s %s: %s\n", LangMap[langType]["commentsTag"], stepIdent, stepExpect)

		codeLine += LangMap[langType]["commentsTag"] + "CODE: 输出验证点实际结果\n"

		*srcCode = append(*srcCode, codeLine)
	}

	if isGroup {
		for _, tsChild := range ts.Steps {
			DealwithTestStep(tsChild, langType, level+1, stepWidth, checkPointIndex, steps, expects, srcCode)
		}
	}
}

package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/biz"
	"github.com/easysoft/zentaoatf/src/biz/zentao"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenFromCmd(url string, entityType string, entityVal string, langType string, singleFile bool,
	account string, password string) {
	params := make(map[string]string)

	params["entityType"] = entityType
	params["entityVal"] = entityVal

	url = utils.UpdateUrl(url)
	zentao.Login(url, account, password)

	var name string
	var testcases []model.TestCase
	if entityType == "product" {
		product := zentao.GetProductInfo(url, params["entityVal"])
		name = product.Name
		testcases = zentao.ListCaseByProduct(url, params["entityVal"])
	} else {
		//taskJson := zentao.GetTaskInfo(url, params["entityVal"])
		//name, _ = taskJson.Get("name").String()
		//json = zentao.ListCaseByProduct(url, params["entityVal"])
	}

	if testcases != nil {
		count, err := Generate(testcases, langType, singleFile, account, password)
		if err == nil {
			utils.SaveConfig("", url, params["entityType"], params["entityVal"], langType, singleFile,
				name, account, password)

			fmt.Sprintf("success to generate %d test scripts in '%s' at %s",
				count, utils.ScriptDir, utils.DateTimeStr(time.Now()))
		} else {
			fmt.Sprintf(err.Error())
		}
	}
}

func Generate(testcases []model.TestCase, langType string, singleFile bool,
	account string, password string) (int, error) {

	casePaths := make([]string, 0)
	for _, cs := range testcases {
		DealwithTestCase(cs, langType, singleFile, &casePaths)
	}

	biz.GenSuite(casePaths)

	return len(testcases), nil
}

func DealwithTestCase(cs model.TestCase, langType string, singleFile bool, casePaths *[]string) {
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

	caseId := cs.Id
	caseTitle := cs.Title

	scriptFile := fmt.Sprintf(utils.ScriptDir+"tc-%s.%s", caseId, LangMap[langType]["extName"])
	if utils.FileExist(scriptFile) {
		scriptFile = fmt.Sprintf(utils.ScriptDir+"tc-%s.%s",
			caseId+"-"+utils.DateTimeStrLong(time.Now()), LangMap[langType]["extName"])
	}

	utils.MkDirIfNeeded(utils.Prefer.WorkDir + utils.ScriptDir)
	*casePaths = append(*casePaths, scriptFile)
	scriptFullPath := utils.Prefer.WorkDir + scriptFile

	steps := make([]string, 0)
	expects := make([]string, 0)
	srcCode := make([]string, 0)

	steps = append(steps, "@开头的为含验证点的步骤")

	temp := fmt.Sprintf("\n%sCODE: 此处编写操作步骤代码\n", LangMap[langType]["commentsTag"])
	srcCode = append(srcCode, temp)

	readme := utils.ReadResData("res/template/readme.tpl") + "\n"

	StepWidth := 20
	stepDisplayMaxWidth := 0
	ComputerTestStepWidth(cs.StepArr, &stepDisplayMaxWidth, StepWidth)

	for _, ts := range cs.StepArr {
		DealwithTestStep(ts, langType, StepWidth, &steps, &expects, &srcCode)
	}

	var expectsTxt string
	if singleFile {
		expectsTxt = strings.Join(expects, "\n")
	} else {
		expectFile := utils.ScriptToExpectName(scriptFullPath)

		expectsTxt = "@file\n"
		utils.WriteFile(expectFile, strings.Join(expects, "\n"))
	}

	path := fmt.Sprintf("res%stemplate%s", string(os.PathSeparator), string(os.PathSeparator))
	template := utils.ReadResData(path + langType + ".tpl")

	id, _ := strconv.Atoi(caseId)
	content := fmt.Sprintf(template,
		id, caseTitle,
		strings.Join(steps, "\n"), expectsTxt,
		readme,
		strings.Join(srcCode, "\n"))

	//fmt.Println(content)

	utils.WriteFile(scriptFullPath, content)
}

func ComputerTestStepWidth(steps []model.TestStep, stepSDisplayMaxWidth *int, stepWidth int) {
	for _, ts := range steps {
		length := len(ts.Id)
		if length > *stepSDisplayMaxWidth {
			*stepSDisplayMaxWidth = length
		}
	}
	*stepSDisplayMaxWidth += stepWidth // prefix space and @step
}

func DealwithTestStep(ts model.TestStep, langType string, stepWidth int,
	steps *[]string, expects *[]string, srcCode *[]string) {
	LangMap := script.GetLangMap()

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
}

package action

import (
	"fmt"
	"github.com/bitly/go-simplejson"
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
	var json *simplejson.Json
	if entityType == "product" {
		productJson := zentao.GetProductInfo(url, params["entityVal"])
		name, _ = productJson.Get("name").String()
		json = zentao.ListCaseByProduct(url, params["entityVal"])
	} else {
		//taskJson := zentao.GetTaskInfo(url, params["entityVal"])
		//name, _ = taskJson.Get("name").String()
		//json = zentao.ListCaseByProduct(url, params["entityVal"])
	}

	if json != nil {
		count, err := Generate(json, url, entityType, entityVal, langType, singleFile, account, password)
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

func Generate(json *simplejson.Json, url string, entityType string, entityVal string, langType string, singleFile bool,
	account string, password string) (int, error) {

	mp, _ := json.Map()
	casePaths := make([]string, 0)
	for _, csJson := range mp {
		if cs, ok := csJson.(map[string]interface{}); ok {
			DealwithTestCase(cs, langType, singleFile, &casePaths)
		}
	}

	biz.GenSuite(casePaths)

	return len(mp), nil

	return 0, nil
}

func DealwithTestCase(tc map[string]interface{}, langType string, singleFile bool, casePaths *[]string) {
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

	//StepWidth := 20

	caseId := tc["id"].(string)
	caseTitle := tc["title"]

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

	//stepDisplayMaxWidth := 0
	//ComputerTestStepWidth(tc["steps], &stepDisplayMaxWidth, StepWidth)
	//
	//level := 1
	//checkPointIndex := 0
	//for _, ts := range tc.Steps {
	//	DealwithTestStep(ts, langType, level, StepWidth, &checkPointIndex, &steps, &expects, &srcCode)
	//}

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
	content := fmt.Sprintf(template,
		caseId, caseTitle,
		strings.Join(steps, "\n"), expectsTxt,
		readme,
		strings.Join(srcCode, "\n"))

	//fmt.Println(content)

	utils.WriteFile(scriptFullPath, content)
}

func ComputerTestStepWidth(steps []model.TestStep, stepSDisplayMaxWidth *int, stepWidth int) {
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

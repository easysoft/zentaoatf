package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/lang"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
	"regexp"
	"strconv"
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

	info := make([]string, 0)
	steps := make([]string, 0)
	independentExpects := make([]string, 0)
	srcCode := make([]string, 0)

	info = append(info, fmt.Sprintf("title=%s", caseTitle))
	info = append(info, fmt.Sprintf("cid=%s", caseId))
	info = append(info, fmt.Sprintf("pid=%s", productId))

	StepWidth := 20
	stepDisplayMaxWidth := 0
	computerTestStepWidth(cs.StepArr, &stepDisplayMaxWidth, StepWidth)

	GenerateTestStepAndScript(cs.StepArr, langType, StepWidth, &steps, &independentExpects, &srcCode)
	info = append(info, strings.Join(steps, "\n"))

	if fileUtils.FileExist(scriptFile) { // update title and steps
		content := fileUtils.ReadFile(scriptFile)

		// replace info
		re, _ := regexp.Compile(`(?s)\[case\].*\[esac\]`)
		content = re.ReplaceAllString(content, "[case]\n"+strings.Join(info, "\n")+"\n\n[esac]")

		fileUtils.WriteFile(scriptFile, content)
		return
	}

	temp := fmt.Sprintf("\n%sCODE: %s", langUtils.LangMap[langType]["commentsTag"],
		i118Utils.I118Prt.Sprintf("your_codes_here"))
	srcCode = append(srcCode, temp)

	//var expectsTxt string
	//if !independentFile {
	//	expectsTxt = strings.Join(independentExpects, "\n")
	//} else {
	//	expectFile := zentaoUtils.ScriptToExpectName(scriptFile)
	//
	//	expectsTxt = "@file\n"
	//	fileUtils.WriteFile(expectFile, strings.Join(independentExpects, "\n"))
	//}

	path := fmt.Sprintf("res%stemplate%s", string(os.PathSeparator), string(os.PathSeparator))
	template := fileUtils.ReadResData(path + langType + ".tpl")

	content := fmt.Sprintf(template, strings.Join(info, "\n"), strings.Join(srcCode, "\n"))

	fileUtils.WriteFile(scriptFile, content)
}

func GenerateTestStepAndScript(teststeps []model.TestStep, langType string, stepWidth int,
	steps *[]string, independentExpects *[]string, srcCode *[]string) {

	var currGroupId string

	groupNo := 0
	childNo := 1
	for idx, ts := range teststeps {
		if idx == 0 { // new group
			groupNo++
			*steps = append(*steps, "")

			if ts.Type == "group" {
				currGroupId = ts.Id
				*steps = append(*steps, fmt.Sprintf("[1. %s]", ts.Desc))
			} else {
				currGroupId = "0"
				*steps = append(*steps, "[group]")
				*steps = append(*steps, zentaoService.GetCaseContent(ts, strconv.Itoa(groupNo))...)
			}

			childNo = 1
			continue
		}

		if ts.Type == "group" { // new group
			groupNo++
			*steps = append(*steps, "")

			currGroupId = ts.Id
			*steps = append(*steps, fmt.Sprintf("[%d. %s]", groupNo, ts.Desc))

			childNo = 1
			continue
		}

		if ts.Type != "group" && ts.Parent != currGroupId { // new group
			groupNo++
			*steps = append(*steps, "")

			currGroupId = "0"
			*steps = append(*steps, "[group]")
			*steps = append(*steps, zentaoService.GetCaseContent(ts, strconv.Itoa(groupNo))...)

			childNo = 1
			continue
		}

		// follow pre group
		var numb string
		if ts.Parent == "0" {
			groupNo++
			numb = fmt.Sprintf("%d", groupNo)
		} else {
			numb = fmt.Sprintf("%d.%d", groupNo, childNo)
		}

		*steps = append(*steps, zentaoService.GetCaseContent(ts, numb)...)
		childNo++
	}

	//isCheckpoint := ts.Expect != ""
	// 处理expects
	//if isCheckpoint {
	//	expectsLine := ""
	//
	//	expectsLine = "# " + stepIdent + " \n"
	//	expectsLine += "CODE: " + i118Utils.I118Prt.Sprintf("expect_result_here")
	//
	//	*independentExpects = append(*independentExpects, expectsLine)
	//}

	// 处理srcCode
	//if isCheckpoint {
	//	codeLine := langUtils.LangMap[langType]["printGrammar"]
	//
	//	codeLine += fmt.Sprintf("  %s %s: %s\n", langUtils.LangMap[langType]["commentsTag"], stepIdent, stepExpect)
	//
	//	codeLine += langUtils.LangMap[langType]["commentsTag"] + "CODE: " + i118Utils.I118Prt.Sprintf("actual_result_here")
	//
	//	*srcCode = append(*srcCode, codeLine)
	//}
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

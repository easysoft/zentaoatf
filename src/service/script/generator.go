package scriptService

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/lang"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
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
	srcCode := fmt.Sprintf("%s %s", langUtils.LangMap[langType]["commentsTag"],
		i118Utils.I118Prt.Sprintf("find_example", string(os.PathSeparator), langType))

	info = append(info, fmt.Sprintf("title=%s", caseTitle))
	info = append(info, fmt.Sprintf("cid=%s", caseId))
	info = append(info, fmt.Sprintf("pid=%s", productId))

	StepWidth := 20
	stepDisplayMaxWidth := 0
	computerTestStepWidth(cs.StepArr, &stepDisplayMaxWidth, StepWidth)

	generateTestStepAndScript(cs.StepArr, &steps, &independentExpects, independentFile)
	info = append(info, strings.Join(steps, "\n"))

	if independentFile {
		expectFile := zentaoUtils.ScriptToExpectName(scriptFile)
		fileUtils.WriteFile(expectFile, strings.Join(independentExpects, "\n"))
	}

	if fileUtils.FileExist(scriptFile) { // update title and steps
		content := fileUtils.ReadFile(scriptFile)

		// replace info
		re, _ := regexp.Compile(`(?s)\[case\].*\[esac\]`)
		content = re.ReplaceAllString(content, "[case]\n"+strings.Join(info, "\n")+"\n\n[esac]")

		fileUtils.WriteFile(scriptFile, content)
		return
	}

	path := fmt.Sprintf("res%stemplate%s", string(os.PathSeparator), string(os.PathSeparator))
	template := fileUtils.ReadResData(path + langType + ".tpl")

	content := fmt.Sprintf(template, strings.Join(info, "\n"), srcCode)

	fileUtils.WriteFile(scriptFile, content)
}

func generateTestStepAndScript(teststeps []model.TestStep, steps *[]string, independentExpects *[]string, independentFile bool) {
	nestedSteps := make([]model.TestStep, 0)
	currGroup := model.TestStep{}
	idx := 0
	for true {
		if idx >= len(teststeps) {
			break
		}

		ts := teststeps[idx]
		if ts.Parent == "0" && ts.Type != "group" { // flat step
			currGroup = model.TestStep{Id: "-1", Desc: "group", Children: make([]model.TestStep, 0)}
			currGroup.Children = append(currGroup.Children, ts)
			idx++

			mutiLine := false
			for true {
				if idx >= len(teststeps) {
					currGroup.MutiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}

				child := teststeps[idx]
				if child.Type != "group" { // flat step
					if !mutiLine {
						mutiLine = zentaoService.IsMutiLine(child)
					}

					currGroup.Children = append(currGroup.Children, child)
				} else { // found a group step
					currGroup.MutiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}
				idx++
			}
		} else if ts.Type == "group" {
			currGroup = model.TestStep{Desc: ts.Desc, Children: make([]model.TestStep, 0)}
			idx++

			mutiLine := false
			for true {
				if idx >= len(teststeps) {
					nestedSteps = append(nestedSteps, currGroup)
					break
				}

				child := teststeps[idx]
				if child.Type != "group" && child.Parent == ts.Id { // child step
					if !mutiLine {
						mutiLine = zentaoService.IsMutiLine(child)
					}

					currGroup.Children = append(currGroup.Children, child)
				} else { // found a group step
					currGroup.MutiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}
				idx++
			}
		}
	}

	stepNumb := 1
	for _, group := range nestedSteps {
		if group.Id == "-1" {
			*steps = append(*steps, "\n[group]")

			for _, child := range group.Children {
				*steps = append(*steps,
					zentaoService.GetCaseContent(child, strconv.Itoa(stepNumb), independentFile, group.MutiLine)...)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, getExcepts(child.Expect))
				}

				stepNumb++
			}
		} else {
			*steps = append(*steps, "\n"+fmt.Sprintf("[%d. %s]", stepNumb, group.Desc))

			for childNo, child := range group.Children {
				numbStr := fmt.Sprintf("%d.%d", stepNumb, childNo+1)
				*steps = append(*steps, zentaoService.GetCaseContent(child, numbStr, independentFile, group.MutiLine)...)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, getExcepts(child.Expect))
				}
			}

			stepNumb++
		}
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

func getExcepts(str string) string {
	str = stringUtils.TrimAll(str)

	arr := strings.Split(str, "\n")

	if len(arr) == 1 {
		return ">> " + str
	} else {
		return ">>\n" + str
	}
}

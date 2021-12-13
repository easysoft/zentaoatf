package scriptUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/lang"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Generate(testcases []model.TestCase, langType string, independentFile bool,
	targetDir string, byModule bool, prefix string) (int, error) {
	caseIds := make([]string, 0)
	for _, cs := range testcases {
		GenerateTestCaseScript(cs, langType, independentFile, &caseIds, targetDir, byModule, prefix)
	}

	GenSuite(caseIds, targetDir)

	return len(testcases), nil
}

func GenerateTestCaseScript(cs model.TestCase, langType string, independentFile bool, caseIds *[]string,
	targetDir string, byModule bool, prefix string) {
	caseId := cs.Id
	productId := cs.Product
	moduleId := cs.Module
	caseTitle := cs.Title

	fileUtils.MkDirIfNeeded(targetDir)
	modulePath := ""
	if byModule && moduleId != "0" {
		modulePath = moduleId + string(os.PathSeparator)
	}

	content := ""
	isOldFormat := false
	scriptFile := fmt.Sprintf(targetDir+"%s%s%s.%s", modulePath, prefix, caseId, langUtils.LangMap[langType]["extName"])
	if fileUtils.FileExist(scriptFile) { // update title and steps
		content = fileUtils.ReadFile(scriptFile)
		isOldFormat = strings.Index(content, "[esac]") > -1
	}

	*caseIds = append(*caseIds, caseId)

	info := make([]string, 0)
	steps := make([]string, 0)
	independentExpects := make([]string, 0)
	srcCode := fmt.Sprintf("%s %s", langUtils.LangMap[langType]["commentsTag"],
		i118Utils.Sprintf("find_example", string(os.PathSeparator), langType))

	info = append(info, fmt.Sprintf("title=%s", caseTitle))
	info = append(info, fmt.Sprintf("cid=%s", caseId))
	info = append(info, fmt.Sprintf("pid=%s", productId))

	StepWidth := 20
	stepDisplayMaxWidth := 0
	computerTestStepWidth(cs.StepArr, &stepDisplayMaxWidth, StepWidth)

	if isOldFormat {
		generateTestStepAndScriptObsolete(cs.StepArr, &steps, &independentExpects, independentFile)
	} else {
		generateTestStepAndScript(cs.StepArr, &steps, &independentExpects, independentFile)
	}
	info = append(info, strings.Join(steps, "\n"))

	if independentFile {
		expectFile := zentaoUtils.ScriptToExpectName(scriptFile)
		fileUtils.WriteFile(expectFile, strings.Join(independentExpects, "\n"))
	}

	if fileUtils.FileExist(scriptFile) { // update title and steps
		regStr := fmt.Sprintf(`(?sm)%s((?U:.*pid.*))\n(.*)%s`,
			constant.LangCommentsRegxMap[langType][0], constant.LangCommentsRegxMap[langType][1])

		// replace info
		re, _ := regexp.Compile(regStr)
		out := re.ReplaceAllString(content, "\n/**\n\n"+
			strings.Join(info, "\n")+"\n\n*/\n")

		fileUtils.WriteFile(scriptFile, out)
		return
	}

	path := fmt.Sprintf("res%stemplate%s", string(os.PathSeparator), string(os.PathSeparator))
	template := fileUtils.ReadResData(path + langType + ".tpl")

	out := fmt.Sprintf(template, strings.Join(info, "\n"), srcCode)
	fileUtils.WriteFile(scriptFile, out)
}

func generateTestStepAndScriptObsolete(testSteps []model.TestStep, steps *[]string, independentExpects *[]string, independentFile bool) {
	nestedSteps := make([]model.TestStep, 0)
	currGroup := model.TestStep{}
	idx := 0

	// convert steps to nested
	for true {
		if idx >= len(testSteps) {
			break
		}

		ts := testSteps[idx]
		if ts.Parent == "0" && ts.Type != "group" { // flat step
			currGroup = model.TestStep{Id: "-1", Desc: "group", Children: make([]model.TestStep, 0)}
			currGroup.Children = append(currGroup.Children, ts)
			idx++

			mutiLine := false
			for true {
				if idx >= len(testSteps) {
					currGroup.MutiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}

				child := testSteps[idx]
				if child.Type != "group" { // flat step
					if !mutiLine {
						mutiLine = zentaoService.IsMultiLine(child)
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
				if idx >= len(testSteps) {
					nestedSteps = append(nestedSteps, currGroup)
					break
				}

				child := testSteps[idx]
				if child.Type != "group" && child.Parent == ts.Id { // child step
					if !mutiLine {
						mutiLine = zentaoService.IsMultiLine(child)
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
	// print nested steps, only one level
	for _, group := range nestedSteps {
		if group.Id == "-1" { // [group]
			*steps = append(*steps, fmt.Sprintf("\n[group]"))

			for _, child := range group.Children {
				*steps = append(*steps,
					zentaoService.GetCaseContent(child, strconv.Itoa(stepNumb), independentFile, group.MutiLine)...)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, getExcepts(child.Expect))
				}

				stepNumb++
			}
		} else { // [1. title]
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

func generateTestStepAndScript(testSteps []model.TestStep, steps *[]string, independentExpects *[]string, independentFile bool) {
	nestedSteps := make([]model.TestStep, 0)

	// convert steps to nested
	for index := 0; index < len(testSteps); index++ {
		ts := testSteps[index]
		item := model.TestStep{Desc: ts.Desc, Expect: ts.Expect, Children: make([]model.TestStep, 0)}

		if ts.Type == "group" {
			nestedSteps = append(nestedSteps, item)
		} else if ts.Type == "item" {
			nestedSteps[len(nestedSteps)-1].Children = append(nestedSteps[len(nestedSteps)-1].Children, item)
		} else if ts.Type == "step" {
			nestedSteps = append(nestedSteps, item)
		}
	}

	// print nested steps, only one level
	stepNumb := 1
	*steps = append(*steps, "")
	for _, item := range nestedSteps {
		numbStr := fmt.Sprintf("%d", stepNumb)
		*steps = append(*steps, zentaoService.GetCaseContent(item, numbStr, independentFile, false)...)

		for childNo, child := range item.Children {
			numbStr := fmt.Sprintf("%d.%d", stepNumb, childNo+1)
			*steps = append(*steps, zentaoService.GetCaseContent(child, numbStr, independentFile, true)...)

			if independentFile && strings.TrimSpace(child.Expect) != "" {
				*independentExpects = append(*independentExpects, getExcepts(child.Expect))
			}
		}

		stepNumb++
	}
}

func GenSuite(cases []string, targetDir string) {
	str := strings.Join(cases, "\n")

	fileUtils.WriteFile(targetDir+"all."+constant.ExtNameSuite, str)
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

package scriptHelper

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	resUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/res"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
	"path/filepath"
	"strconv"
	"strings"
)

func GenerateScripts(cases []commDomain.ZtfCase, langType string, independentFile bool,
	byModule bool, targetDir string) (int, error) {
	caseIds := make([]string, 0)

	if commConsts.ComeFrom == "cmd" { // from cmd
		targetDir = stdinUtils.GetInput("", targetDir, "where_to_store_script", targetDir)
		stdinUtils.InputForBool(&byModule, byModule, "co_organize_by_module")
	}
	targetDir = fileUtils.AbsolutePath(targetDir)

	for _, cs := range cases {
		GenerateScript(cs, langType, independentFile, &caseIds, targetDir, byModule)
	}

	GenSuite(caseIds, targetDir)

	return len(cases), nil
}

func GenerateScript(cs commDomain.ZtfCase, langType string, independentFile bool, caseIds *[]string,
	targetDir string, byModule bool) {
	caseId := cs.Id
	productId := cs.Product
	moduleId := cs.Module
	caseTitle := cs.Title

	if byModule {
		targetDir = filepath.Join(targetDir, strconv.Itoa(moduleId))
	}

	fileUtils.MkDirIfNeeded(targetDir)

	content := ""
	isOldFormat := false
	scriptFile := filepath.Join(targetDir, fmt.Sprintf("%d.%s", caseId, langUtils.LangMap[langType]["extName"]))
	if fileUtils.FileExist(scriptFile) { // update title and steps
		content = fileUtils.ReadFile(scriptFile)
		isOldFormat = strings.Index(content, "[esac]") > -1
	}

	*caseIds = append(*caseIds, strconv.Itoa(caseId))

	info := make([]string, 0)
	steps := make([]string, 0)
	independentExpects := make([]string, 0)
	srcCode := fmt.Sprintf("%s %s", langUtils.LangMap[langType]["commentsTag"],
		i118Utils.Sprintf("find_example", consts.PthSep, langType))

	info = append(info, fmt.Sprintf("title=%s", caseTitle))
	info = append(info, fmt.Sprintf("cid=%s", caseId))
	info = append(info, fmt.Sprintf("pid=%s", productId))

	StepWidth := 20
	stepDisplayMaxWidth := 0
	computerTestStepWidth(cs.Steps, &stepDisplayMaxWidth, StepWidth)

	if isOldFormat {
		generateTestStepAndScriptObsolete(cs.Steps, &steps, &independentExpects, independentFile)
	} else {
		generateTestStepAndScript(cs.Steps, &steps, &independentExpects, independentFile)
	}
	info = append(info, strings.Join(steps, "\n"))

	if independentFile {
		expectFile := ScriptToExpectName(scriptFile)
		fileUtils.WriteFile(expectFile, strings.Join(independentExpects, "\n"))
	}

	if fileUtils.FileExist(scriptFile) { // update title and steps
		newContent := strings.Join(info, "\n")
		ReplaceCaseDesc(newContent, scriptFile)
		return
	}

	path := fmt.Sprintf("res%stemplate%s", consts.PthSep, consts.PthSep)
	template, _ := resUtils.ReadRes(path + langType + ".tpl")

	out := fmt.Sprintf(string(template), strings.Join(info, "\n"), srcCode)
	fileUtils.WriteFile(scriptFile, out)
}

func generateTestStepAndScriptObsolete(testSteps []commDomain.ZtfStep, steps *[]string, independentExpects *[]string, independentFile bool) {
	nestedSteps := make([]commDomain.ZtfStep, 0)
	currGroup := commDomain.ZtfStep{}
	idx := 0

	// convert steps to nested
	for true {
		if idx >= len(testSteps) {
			break
		}

		ts := testSteps[idx]
		if ts.Parent == "0" && ts.Type != "group" { // flat step
			currGroup = commDomain.ZtfStep{Id: "-1", Desc: "group", Children: make([]commDomain.ZtfStep, 0)}
			currGroup.Children = append(currGroup.Children, ts)
			idx++

			mutiLine := false
			for true {
				if idx >= len(testSteps) {
					currGroup.MultiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}

				child := testSteps[idx]
				if child.Type != "group" { // flat step
					if !mutiLine {
						mutiLine = IsMultiLine(child)
					}

					currGroup.Children = append(currGroup.Children, child)
				} else { // found a group step
					currGroup.MultiLine = mutiLine
					nestedSteps = append(nestedSteps, currGroup)
					break
				}
				idx++
			}
		} else if ts.Type == "group" {
			currGroup = commDomain.ZtfStep{Desc: ts.Desc, Children: make([]commDomain.ZtfStep, 0)}
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
						mutiLine = IsMultiLine(child)
					}

					currGroup.Children = append(currGroup.Children, child)
				} else { // found a group step
					currGroup.MultiLine = mutiLine
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
				stepContent, expectContent := GetCaseContent(child, strconv.Itoa(stepNumb), independentFile, group.MultiLine)
				*steps = append(*steps, stepContent)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, expectContent)
				}

				stepNumb++
			}
		} else { // [1. title]
			*steps = append(*steps, "\n"+fmt.Sprintf("[%d. %s]", stepNumb, group.Desc))

			for childNo, child := range group.Children {
				numbStr := fmt.Sprintf("%d.%d", stepNumb, childNo+1)

				stepContent, expectContent := GetCaseContent(child, numbStr, independentFile, group.MultiLine)
				*steps = append(*steps, stepContent)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, expectContent)
				}
			}

			stepNumb++
		}
	}
}

func generateTestStepAndScript(testSteps []commDomain.ZtfStep, steps *[]string, independentExpects *[]string, independentFile bool) {
	nestedSteps := make([]commDomain.ZtfStep, 0)

	// convert steps to nested
	for index := 0; index < len(testSteps); index++ {
		ts := testSteps[index]
		item := commDomain.ZtfStep{Desc: ts.Desc, Expect: ts.Expect, Children: make([]commDomain.ZtfStep, 0)}

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
		stepLines1, expects1 := GetCaseContent(item, numbStr, independentFile, false)
		*steps = append(*steps, stepLines1)

		if independentFile && strings.TrimSpace(item.Expect) != "" {
			*independentExpects = append(*independentExpects, expects1)
		}

		for childNo, child := range item.Children {
			numbStr := fmt.Sprintf("%d.%d", stepNumb, childNo+1)
			stepLines2, expects2 := GetCaseContent(child, numbStr, independentFile, true)
			*steps = append(*steps, stepLines2)

			if independentFile && strings.TrimSpace(child.Expect) != "" {
				*independentExpects = append(*independentExpects, expects2)
			}
		}

		stepNumb++
	}
}

func computerTestStepWidth(steps []commDomain.ZtfStep, stepSDisplayMaxWidth *int, stepWidth int) {
	for _, ts := range steps {
		length := len(ts.Id)
		if length > *stepSDisplayMaxWidth {
			*stepSDisplayMaxWidth = length
		}
	}
	*stepSDisplayMaxWidth += stepWidth // prefix space and @step
}

func GenSuite(cases []string, targetDir string) {
	str := strings.Join(cases, "\n")

	fileUtils.WriteFile(targetDir+"all."+commConsts.ExtNameSuite, str)
}

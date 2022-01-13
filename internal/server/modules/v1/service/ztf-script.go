package service

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	resUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/res"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"regexp"

	"strconv"
	"strings"
)

type ZtfScriptService struct {
	ZtfCaseService *ZtfCaseService   `inject:""`
	ProjectRepo    *repo.ProjectRepo `inject:""`
}

func NewZtfScriptService() *ZtfScriptService {
	return &ZtfScriptService{}
}

func (s *ZtfScriptService) Generate(cases []commDomain.ZtfCase, langType string, independentFile bool,
	byModule bool, targetDir string, prefix string) (int, error) {
	caseIds := make([]string, 0)
	for _, cs := range cases {
		s.GenerateScript(cs, langType, independentFile, &caseIds, targetDir, byModule, prefix)
	}

	s.GenSuite(caseIds, targetDir)

	return len(cases), nil
}

func (s *ZtfScriptService) GenerateScript(cs commDomain.ZtfCase, langType string, independentFile bool, caseIds *[]string,
	targetDir string, byModule bool, prefix string) {
	caseId := cs.Id
	productId := cs.Product
	moduleId := cs.Module
	caseTitle := cs.Title

	fileUtils.MkDirIfNeeded(targetDir)
	modulePath := ""
	if byModule && moduleId != "0" {
		modulePath = fmt.Sprintf("%d%s", moduleId, consts.PthSep)
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
		i118Utils.Sprintf("find_example", consts.PthSep, langType))

	info = append(info, fmt.Sprintf("title=%s", caseTitle))
	info = append(info, fmt.Sprintf("cid=%s", caseId))
	info = append(info, fmt.Sprintf("pid=%s", productId))

	StepWidth := 20
	stepDisplayMaxWidth := 0
	s.computerTestStepWidth(cs.StepArr, &stepDisplayMaxWidth, StepWidth)

	if isOldFormat {
		s.generateTestStepAndScriptObsolete(cs.StepArr, &steps, &independentExpects, independentFile)
	} else {
		s.generateTestStepAndScript(cs.StepArr, &steps, &independentExpects, independentFile)
	}
	info = append(info, strings.Join(steps, "\n"))

	if independentFile {
		expectFile := zentaoUtils.ScriptToExpectName(scriptFile)
		fileUtils.WriteFile(expectFile, strings.Join(independentExpects, "\n"))
	}

	if fileUtils.FileExist(scriptFile) { // update title and steps
		regStr := fmt.Sprintf(`(?sm)%s((?U:.*pid.*))\n(.*)%s`,
			langUtils.LangCommentsRegxMap[langType][0], langUtils.LangCommentsRegxMap[langType][1])

		// replace info
		re, _ := regexp.Compile(regStr)
		newContent := fmt.Sprintf("\n%s\n\n%s\n\n%s\n",
			langUtils.LangCommentsTagMap[langType][0],
			strings.Join(info, "\n"),
			langUtils.LangCommentsTagMap[langType][1])

		out := re.ReplaceAllString(content, newContent)

		fileUtils.WriteFile(scriptFile, out)
		return
	}

	path := fmt.Sprintf("res%stemplate%s", consts.PthSep, consts.PthSep)
	template, _ := resUtils.ReadRes(path + langType + ".tpl")

	out := fmt.Sprintf(string(template), strings.Join(info, "\n"), srcCode)
	fileUtils.WriteFile(scriptFile, out)
}

func (s *ZtfScriptService) generateTestStepAndScriptObsolete(testSteps []commDomain.ZtfStep, steps *[]string, independentExpects *[]string, independentFile bool) {
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
						mutiLine = s.ZtfCaseService.IsMultiLine(child)
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
						mutiLine = s.ZtfCaseService.IsMultiLine(child)
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
				*steps = append(*steps,
					s.ZtfCaseService.GetCaseContent(child, strconv.Itoa(stepNumb), independentFile, group.MultiLine)...)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, s.getExcepts(child.Expect))
				}

				stepNumb++
			}
		} else { // [1. title]
			*steps = append(*steps, "\n"+fmt.Sprintf("[%d. %s]", stepNumb, group.Desc))

			for childNo, child := range group.Children {
				numbStr := fmt.Sprintf("%d.%d", stepNumb, childNo+1)
				*steps = append(*steps, s.ZtfCaseService.GetCaseContent(child, numbStr, independentFile, group.MultiLine)...)

				if independentFile && strings.TrimSpace(child.Expect) != "" {
					*independentExpects = append(*independentExpects, s.getExcepts(child.Expect))
				}
			}

			stepNumb++
		}
	}
}

func (s *ZtfScriptService) generateTestStepAndScript(testSteps []commDomain.ZtfStep, steps *[]string, independentExpects *[]string, independentFile bool) {
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
		*steps = append(*steps, s.ZtfCaseService.GetCaseContent(item, numbStr, independentFile, false)...)

		for childNo, child := range item.Children {
			numbStr := fmt.Sprintf("%d.%d", stepNumb, childNo+1)
			*steps = append(*steps, s.ZtfCaseService.GetCaseContent(child, numbStr, independentFile, true)...)

			if independentFile && strings.TrimSpace(child.Expect) != "" {
				*independentExpects = append(*independentExpects, s.getExcepts(child.Expect))
			}
		}

		stepNumb++
	}
}

func (s *ZtfScriptService) GenSuite(cases []string, targetDir string) {
	str := strings.Join(cases, "\n")

	fileUtils.WriteFile(targetDir+"all."+commConsts.ExtNameSuite, str)
}

func (s *ZtfScriptService) computerTestStepWidth(steps []commDomain.ZtfStep, stepSDisplayMaxWidth *int, stepWidth int) {
	for _, ts := range steps {
		length := len(ts.Id)
		if length > *stepSDisplayMaxWidth {
			*stepSDisplayMaxWidth = length
		}
	}
	*stepSDisplayMaxWidth += stepWidth // prefix space and @step
}

func (s *ZtfScriptService) getExcepts(str string) string {
	str = stringUtils.TrimAll(str)

	arr := strings.Split(str, "\n")

	if len(arr) == 1 {
		return ">> " + str
	} else {
		return ">>\n" + str
	}
}

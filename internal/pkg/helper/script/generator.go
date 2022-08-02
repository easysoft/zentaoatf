package scriptHelper

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"path/filepath"
	"strconv"
	"strings"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	"github.com/easysoft/zentaoatf/pkg/consts"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	resUtils "github.com/easysoft/zentaoatf/pkg/lib/res"
	stdinUtils "github.com/easysoft/zentaoatf/pkg/lib/stdin"
)

func GenerateScripts(cases []commDomain.ZtfCase, langType string, independentFile bool,
	byModule bool, targetDir string) (pths []string, realPath string, err error) {
	caseIds := make([]string, 0)

	if commConsts.ExecFrom == commConsts.FromCmd { // from cmd
		targetDir = stdinUtils.GetInput("", targetDir, "where_to_store_script", targetDir)
		stdinUtils.InputForBool(&byModule, byModule, "co_organize_by_module")
	}
	targetDir = fileUtils.AbsolutePath(targetDir)
	realPath = targetDir

	createNew := false
	for _, cs := range cases {
		pth, _ := GenerateScript(cs, langType, independentFile, &caseIds, targetDir, byModule)
		pths = append(pths, pth)

		if cs.ScriptPath == "" {
			createNew = true
		}
	}

	if createNew {
		GenSuite(caseIds, targetDir)
	}

	return
}

func GenerateScript(cs commDomain.ZtfCase, langType string, independentFile bool, caseIds *[]string,
	targetDir string, byModule bool) (scriptPath string, err error) {

	caseId := cs.Id
	productId := cs.Product
	moduleId := cs.Module
	caseTitle := cs.Title

	if byModule {
		targetDir = filepath.Join(targetDir, strconv.Itoa(moduleId))
	}

	scriptPath = cs.ScriptPath
	if scriptPath == "" {
		fileUtils.MkDirIfNeeded(targetDir)
		scriptPath = filepath.Join(targetDir, fmt.Sprintf("%d.%s", caseId, commConsts.LangMap[langType]["extName"]))
	}

	*caseIds = append(*caseIds, strconv.Itoa(caseId))

	info := make([]string, 0)
	steps := make([]string, 0)
	independentExpects := make([]string, 0)
	srcCode := fmt.Sprintf("%s %s", commConsts.LangMap[langType]["commentsTag"],
		i118Utils.Sprintf("find_example", consts.FilePthSep, langType))

	info = append(info, fmt.Sprintf("title=%s", caseTitle))
	info = append(info, fmt.Sprintf("cid=%d", caseId))
	info = append(info, fmt.Sprintf("pid=%d", productId))

	StepWidth := 20
	stepDisplayMaxWidth := 0
	computerTestStepWidth(cs.Steps, &stepDisplayMaxWidth, StepWidth)
	generateTestStepAndScript(cs.Steps, &steps, &independentExpects, independentFile)

	info = append(info, strings.Join(steps, "\n"))

	if independentFile {
		expectFile := ScriptToExpectName(scriptPath)
		fileUtils.WriteFile(expectFile, strings.Join(independentExpects, "\n"))
	}

	if fileUtils.FileExist(scriptPath) { // update title and steps
		newContent := strings.Join(info, "\n")
		ReplaceCaseDesc(newContent, scriptPath)
		return
	}

	templatePath := fmt.Sprintf("res%stemplate%s", consts.FilePthSep, consts.FilePthSep)
	template, _ := resUtils.ReadRes(templatePath + langType + ".tpl")

	out := fmt.Sprintf(string(template), strings.Join(info, "\n"), srcCode)
	fileUtils.WriteFile(scriptPath, out)

	return
}

func GenEmptyScript(name, lang, pth string, productId int) {
	srcCode := fmt.Sprintf("%s %s", commConsts.LangMap[lang]["commentsTag"],
		i118Utils.Sprintf("find_example", consts.FilePthSep, lang))

	info := make([]string, 0)
	info = append(info, fmt.Sprintf("title=%s", name))
	info = append(info, fmt.Sprintf("cid=%d", 0))
	info = append(info, fmt.Sprintf("pid=%d", productId))

	templatePath := fmt.Sprintf("res%stemplate%s", consts.FilePthSep, consts.FilePthSep)
	template, _ := resUtils.ReadRes(templatePath + lang + ".tpl")

	out := fmt.Sprintf(string(template), strings.Join(info, "\n"), srcCode)
	fileUtils.WriteFile(pth, out)
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
		length := len(strconv.Itoa(ts.Id))
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

package zentaoHelper

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	"github.com/fatih/color"
)

func Checkout(settings commDomain.SyncSettings, config commDomain.WorkspaceConf, workspacePath string) (
	pths []string, err error) {

	productId := settings.ProductId
	moduleId := settings.ModuleId
	suiteId := settings.SuiteId
	taskId := settings.TaskId
	caseId := settings.CaseId
	casePath := settings.CasePath

	byModule := settings.SaveByModule
	independentFile := settings.IndependentFile
	lang := settings.Lang

	ok := langHelper.CheckSupportLanguages(lang)
	if !ok {
		return
	}
	spew.Dump(settings)
	cases := make([]commDomain.ZtfCase, 0)
	if caseId != 0 {
		cs, err := GetTestCaseDetail(caseId, config)
		cs.ScriptPath = casePath
		if err == nil {
			cases = append(cases, cs)
		}
	} else if len(settings.CaseIds) > 0 {
		cases, err = LoadTestCasesDetailByCaseIds(settings.CaseIds, config)
	} else {
		cases, err = LoadTestCasesDetail(productId, moduleId, suiteId, taskId, config)
	}

	if err != nil {
		return
	}

	if len(cases) == 0 {
		err = errors.New("no_cases_found")
		return
	}

	if productId == 0 {
		productId = cases[0].Product
	}
	targetDir := fileUtils.AddFilePathSepIfNeeded(filepath.Join(workspacePath, fmt.Sprintf("product%d", productId)))

	pths, targetDir, err = scriptHelper.GenerateScripts(cases, lang, independentFile, byModule, targetDir)
	if err == nil {
		logUtils.Infof(i118Utils.Sprintf("success_to_generate", len(pths), targetDir))
	} else {
		logUtils.Infof(color.RedString(err.Error()))
	}

	return
}

func CheckIn(cases []string, config commDomain.WorkspaceConf, noNeedConfirm, withCode bool) (count int, err error) {
	for _, cs := range cases {
		pass, id, _, title, _ := scriptHelper.GetCaseInfo(cs)
		if !pass {
			continue
		}

		steps := scriptHelper.GetStepAndExpectMap(cs)
		script, _ := scriptHelper.GetScriptContent(cs, -1)
		err = CommitCase(id, title, steps, script, config, noNeedConfirm, withCode)

		if err == nil {
			count++
		}
	}

	logUtils.Infof(i118Utils.Sprintf("commit_cases_result", count))

	return
}

package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	stdinUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/stdin"
)

type TestScriptService struct {
}

func NewTestScriptService() *TestScriptService {
	return &TestScriptService{}
}

func (s *TestScriptService) GenerateScripts(cases []commDomain.ZtfCase, langType string, independentFile bool,
	byModule bool, targetDir string) (int, error) {
	caseIds := make([]string, 0)

	if commConsts.ComeFrom == "cmd" { // from cmd
		targetDir = stdinUtils.GetInput("", targetDir, "where_to_store_script", targetDir)
		stdinUtils.InputForBool(&byModule, byModule, "co_organize_by_module")
	}
	targetDir = fileUtils.AbsolutePath(targetDir)

	for _, cs := range cases {
		scriptUtils.GenerateScript(cs, langType, independentFile, &caseIds, targetDir, byModule)
	}

	scriptUtils.GenSuite(caseIds, targetDir)

	return len(cases), nil
}

package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
)

type TestScriptService struct {
}

func NewTestScriptService() *TestScriptService {
	return &TestScriptService{}
}

func (s *TestScriptService) GenerateScripts(cases []commDomain.ZtfCase, langType string, independentFile bool,
	byModule bool, targetDir string, prefix string) (int, error) {
	caseIds := make([]string, 0)

	for _, cs := range cases {
		scriptUtils.GenerateScript(cs, langType, independentFile, &caseIds, targetDir, byModule, prefix)
	}

	scriptUtils.GenSuite(caseIds, targetDir)

	return len(cases), nil
}

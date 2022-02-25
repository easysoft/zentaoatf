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

func (s *TestScriptService) GenerateScripts(cases []commDomain.ZtfCase, langType string, independentFile,
	byModule bool, targetDir string) (int, error) {

	return scriptUtils.GenerateScripts(cases, langType, independentFile, byModule, targetDir)
}

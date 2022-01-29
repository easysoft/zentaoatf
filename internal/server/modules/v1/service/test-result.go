package service

import (
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	analysisUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/analysis"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/zentao"
)

type TestResultService struct {
}

func NewTestResultService() *TestResultService {
	return &TestResultService{}
}

func (s *TestResultService) Submit(result serverDomain.ZentaoResult, projectPath string) (err error) {
	report, err := analysisUtils.ReadReport(projectPath, result.Seq)
	if err != nil {
		return
	}

	zentaoUtils.CommitResult(report, result, projectPath)

	return
}

package service

import (
	analysisUtils "github.com/aaronchen2k/deeptest/internal/server/modules/utils/analysis"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/utils/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
)

type TestResultService struct {
}

func NewTestResultService() *TestResultService {
	return &TestResultService{}
}

func (s *TestResultService) Submit(result serverDomain.ZentaoResultSubmitReq, projectPath string) (err error) {
	report, err := analysisUtils.ReadReport(projectPath, result.Seq)
	if err != nil {
		return
	}

	err = zentaoUtils.CommitResult(report, result.ProductId, result.TaskId, projectPath)

	return
}

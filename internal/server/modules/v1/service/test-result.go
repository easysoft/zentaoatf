package service

import (
	analysisUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/analysis"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
)

type TestResultService struct {
}

func NewTestResultService() *TestResultService {
	return &TestResultService{}
}

func (s *TestResultService) Submit(result serverDomain.ZentaoResultSubmitReq, workspacePath string) (err error) {
	report, err := analysisUtils.ReadReportByWorkspaceSeq(workspacePath, result.Seq)
	if err != nil {
		return
	}

	err = zentaoUtils.CommitResult(report, result.ProductId, result.TaskId, workspacePath)

	return
}

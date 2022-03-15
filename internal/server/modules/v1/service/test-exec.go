package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	analysisUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/analysis"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/jinzhu/copier"
	"path/filepath"
)

type TestExecService struct {
	WorkspaceRepo *repo.WorkspaceRepo `inject:""`
}

func NewTestExecService() *TestExecService {
	return &TestExecService{}
}

func (s *TestExecService) List(workspacePath string) (ret []serverDomain.TestReportSummary, err error) {
	reportFiles := analysisUtils.ListReport(workspacePath)

	for _, seq := range reportFiles {
		var summary serverDomain.TestReportSummary

		report, err1 := analysisUtils.ReadReportByWorkspaceSeq(workspacePath, seq)
		if err1 != nil { // ignore wrong json result
			continue
		}
		copier.Copy(&summary, report)

		summary.Seq = seq
		ret = append(ret, summary)
	}

	return
}

func (s *TestExecService) Get(workspacePath string, seq string) (report commDomain.ZtfReport, err error) {
	return analysisUtils.ReadReportByWorkspaceSeq(workspacePath, seq)
}

func (s *TestExecService) Delete(workspacePath string, seq string) (err error) {
	dir := filepath.Join(workspacePath, commConsts.LogDirName)

	di := filepath.Join(dir, seq)
	err = fileUtils.RmDir(di)

	return
}

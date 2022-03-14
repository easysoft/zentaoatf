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

type SiteService struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func NewSiteService() *SiteService {
	return &SiteService{}
}

func (s *SiteService) List(projectPath string) (ret []serverDomain.TestReportSummary, err error) {
	reportFiles := analysisUtils.ListReport(projectPath)

	for _, seq := range reportFiles {
		var summary serverDomain.TestReportSummary

		report, err1 := analysisUtils.ReadReportByProjectSeq(projectPath, seq)
		if err1 != nil { // ignore wrong json result
			continue
		}
		copier.Copy(&summary, report)

		summary.Seq = seq
		ret = append(ret, summary)
	}

	return
}

func (s *SiteService) Get(projectPath string, seq string) (report commDomain.ZtfReport, err error) {
	return analysisUtils.ReadReportByProjectSeq(projectPath, seq)
}

func (s *SiteService) Delete(projectPath string, seq string) (err error) {
	dir := filepath.Join(projectPath, commConsts.LogDirName)

	di := filepath.Join(dir, seq)
	err = fileUtils.RmDir(di)

	return
}

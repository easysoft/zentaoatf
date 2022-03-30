package service

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	analysisUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/analysis"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
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

func (s *TestExecService) Paginate(siteId, productId int, req serverDomain.ReqPaginate) (
	data domain.PageData, err error) {

	reports := []serverDomain.TestReportSummary{}

	workspaces, _ := s.WorkspaceRepo.ListWorkspacesByProduct(siteId, productId)

	pageNo := req.Page - 1
	pageSize := req.PageSize
	jumpNo := pageNo * pageSize

	count := 0
	for _, workspace := range workspaces {
		if workspace.Type != commConsts.ZTF {
			continue
		}

		reportFiles := analysisUtils.ListReport(workspace.Path)
		for _, seq := range reportFiles {
			if count < jumpNo || len(reports) >= pageSize {
				count += 1
				continue
			}

			summary := serverDomain.TestReportSummary{WorkspaceId: int(workspace.ID)}

			report, err1 := analysisUtils.ReadReportByWorkspaceSeq(workspace.Path, seq)
			if err1 != nil { // ignore wrong json result
				continue
			}
			copier.Copy(&summary, report)

			summary.No = fmt.Sprintf("%d-%s", workspace.ID, seq)
			summary.Seq = seq
			summary.WorkspaceName = workspace.Name
			reports = append(reports, summary)

			count += 1
		}
	}

	data.Populate(reports, int64(count), req.Page, req.PageSize)

	return
}

func (s *TestExecService) Get(workspaceId int, seq string) (report commDomain.ZtfReport, err error) {
	workspace, _ := s.WorkspaceRepo.FindById(uint(workspaceId))

	return analysisUtils.ReadReportByWorkspaceSeq(workspace.Path, seq)
}

func (s *TestExecService) Delete(workspaceId int, seq string) (err error) {
	workspace, _ := s.WorkspaceRepo.FindById(uint(workspaceId))
	dir := filepath.Join(workspace.Path, commConsts.LogDirName)

	di := filepath.Join(dir, seq)
	err = fileUtils.RmDir(di)

	return
}

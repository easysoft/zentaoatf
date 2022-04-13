package service

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/comm/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/pkg/domain"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	"github.com/jinzhu/copier"
	"path/filepath"
)

type TestResultService struct {
	SiteRepo      *repo.SiteRepo      `inject:""`
	WorkspaceRepo *repo.WorkspaceRepo `inject:""`
}

func NewTestResultService() *TestResultService {
	return &TestResultService{}
}

func (s *TestResultService) Paginate(siteId, productId uint, req serverDomain.ReqPaginate) (
	data domain.PageData, err error) {

	reports := []serverDomain.TestReportSummary{}

	workspaces, _ := s.WorkspaceRepo.ListByProduct(siteId, productId)

	pageNo := req.Page - 1
	pageSize := req.PageSize
	jumpNo := pageNo * pageSize

	count := 0
	for _, workspace := range workspaces {
		//if workspace.Type != commConsts.ZTF {
		//	continue
		//}

		reportFiles := analysisHelper.ListReport(workspace.Path)
		for _, seq := range reportFiles {
			if count < jumpNo || len(reports) >= pageSize {
				count += 1
				continue
			}

			summary := serverDomain.TestReportSummary{WorkspaceId: int(workspace.ID)}

			report, err1 := analysisHelper.ReadReportByWorkspaceSeq(workspace.Path, seq)
			if err1 != nil { // ignore wrong json result
				continue
			}
			copier.Copy(&summary, report)

			summary.No = fmt.Sprintf("%d-%s", workspace.ID, seq)
			summary.Seq = seq
			summary.WorkspaceId = int(workspace.ID)
			summary.WorkspaceName = workspace.Name
			reports = append(reports, summary)

			count += 1
		}
	}

	data.Populate(reports, int64(count), req.Page, req.PageSize)

	return
}

func (s *TestResultService) Get(workspaceId int, seq string) (report commDomain.ZtfReport, err error) {
	workspace, _ := s.WorkspaceRepo.Get(uint(workspaceId))
	report, err = analysisHelper.ReadReportByWorkspaceSeq(workspace.Path, seq)
	report.WorkspaceId = workspaceId
	report.Seq = seq

	return
}

func (s *TestResultService) Delete(workspaceId int, seq string) (err error) {
	workspace, _ := s.WorkspaceRepo.Get(uint(workspaceId))
	dir := filepath.Join(workspace.Path, commConsts.LogDirName)

	di := filepath.Join(dir, seq)
	err = fileUtils.RmDir(di)

	return
}

func (s *TestResultService) Submit(result serverDomain.ZentaoResultSubmitReq, siteId, productId int) (err error) {
	site, err := s.SiteRepo.Get(uint(siteId))

	workspace, _ := s.WorkspaceRepo.Get(uint(result.WorkspaceId))

	report, err := analysisHelper.ReadReportByWorkspaceSeq(workspace.Path, result.Seq)
	if err != nil {
		return
	}

	config := configHelper.LoadBySite(site)
	err = zentaoHelper.CommitResult(report, result.ProductId, result.TaskId, config)

	return
}

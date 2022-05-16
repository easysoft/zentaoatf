package service

import (
	"fmt"
	"io/fs"
	"path/filepath"

	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/comm/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/comm/helper/zentao"
	"github.com/easysoft/zentaoatf/internal/pkg/domain"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	"github.com/jinzhu/copier"
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

		logDir := filepath.Join(workspace.Path, commConsts.LogDirName)
		reportFiles := analysisHelper.ListReportByModTime(logDir)
		for _, fi := range reportFiles {
			if count < jumpNo || len(reports) >= pageSize {
				count += 1
				continue
			}

			seq := fi.Name()
			summary := serverDomain.TestReportSummary{WorkspaceId: int(workspace.ID)}

			report, err1 := analysisHelper.ReadReportByWorkspaceSeq2(workspace.Path, seq, false)
			if err1 != nil { // ignore wrong json result
				continue
			}
			copier.Copy(&summary, report)

			summary.No = fmt.Sprintf("%d-%s", workspace.ID, seq)
			summary.Seq = seq
			summary.WorkspaceId = int(workspace.ID)
			summary.WorkspaceName = workspace.Name

			if report.Total == 1 {
				_, summary.TestScriptName = filepath.Split(report.FuncResult[0].Path)
			}

			reports = append(reports, summary)

			count += 1
		}

		// if the num of current log report is not enough, list bak reports.
		reportLen := len(reports)
		jumpNo += reportLen

		logBakDir := filepath.Join(workspace.Path, commConsts.LogBakDirName)
		bakReportFiles := analysisHelper.ListReportByModTime(logBakDir)
		for _, fi := range bakReportFiles {
			if count < jumpNo || len(reports) >= pageSize {
				count += 1
				continue
			}

			seq := fi.Name()
			summary := serverDomain.TestReportSummary{WorkspaceId: int(workspace.ID)}

			report, err1 := analysisHelper.ReadReportByWorkspaceSeq2(workspace.Path, seq, true)
			if err1 != nil { // ignore wrong json result
				continue
			}
			copier.Copy(&summary, report)

			summary.No = fmt.Sprintf("%d-%s", workspace.ID, seq)
			summary.Seq = seq
			summary.WorkspaceId = int(workspace.ID)
			summary.WorkspaceName = workspace.Name

			if report.Total == 1 {
				_, summary.TestScriptName = filepath.Split(report.FuncResult[0].Path)
			}

			reports = append(reports, summary)

			count += 1
		}
	}

	data.Populate(reports, int64(count), req.Page, req.PageSize)

	return
}

func (s *TestResultService) GetLatest(siteId, productId uint) (summary serverDomain.TestReportSummary, err error) {
	workspaces, _ := s.WorkspaceRepo.ListByProduct(siteId, productId)

	var fi fs.FileInfo
	var ws model.Workspace
	for _, workspace := range workspaces {
		reportFiles := analysisHelper.ListReportByModTime(workspace.Path + commConsts.LogDir)
		if len(reportFiles) == 0 {
			continue
		}

		if fi != nil && fi.ModTime().Before(reportFiles[0].ModTime()) {
			continue
		}

		fi = reportFiles[0]
		ws = workspace
	}

	summary = serverDomain.TestReportSummary{WorkspaceId: int(ws.ID)}
	seq := fi.Name()
	report, err1 := analysisHelper.ReadReportByWorkspaceSeq(ws.Path, seq)
	if err1 != nil { // ignore wrong json result
		return
	}

	copier.Copy(&summary, report)
	summary.No = fmt.Sprintf("%d-%s", ws.ID, seq)
	summary.Seq = seq
	summary.WorkspaceId = int(ws.ID)
	summary.WorkspaceName = ws.Name

	if report.Total == 1 {
		_, summary.TestScriptName = filepath.Split(report.FuncResult[0].Path)
	}

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
	err = zentaoHelper.CommitResult(report, result.ProductId, result.TaskId, config, nil)

	return
}

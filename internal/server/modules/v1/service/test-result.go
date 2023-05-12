package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/pkg/domain"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/jinzhu/copier"
)

type TestResultService struct {
	SiteRepo         *repo.SiteRepo      `inject:""`
	WorkspaceRepo    *repo.WorkspaceRepo `inject:""`
	ProxyService     *ProxyService       `inject:""`
	WorkspaceService *WorkspaceService   `inject:""`
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

	maxSize := (pageNo + 1) * pageSize
	for _, workspace := range workspaces {
		reportSeqs := analysisHelper.ListReport2(workspace.Path, maxSize)
		for _, seq := range reportSeqs {
			summary := serverDomain.TestReportSummary{
				WorkspaceId: int(workspace.ID),
			}

			report, _, err1 := analysisHelper.ReadReportByWorkspaceSeq(workspace.Path, seq)
			if err1 != nil { // ignore wrong json result
				continue
			}
			copier.Copy(&summary, report)

			summary.No = fmt.Sprintf("%d-%s", workspace.ID, seq)
			summary.Seq = seq
			summary.WorkspaceId = int(workspace.ID)
			summary.WorkspaceName = workspace.Name

			if report.Total == 1 && len(report.FuncResult) > 0 {
				scriptName := ""
				if commConsts.PthSep == "\\" {
					scriptName = strings.Replace(report.FuncResult[0].Path, "/", "\\", -1)
				} else {
					scriptName = strings.Replace(report.FuncResult[0].Path, "\\", "/", -1)
				}
				_, summary.TestScriptName = filepath.Split(scriptName)
			}

			reports = append(reports, summary)
		}
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].StartTime > reports[j].StartTime
	})

	count := len(reports) - jumpNo
	if count > pageSize {
		count = pageSize
	} else if count < 1 {
		count = 0
	}

	data.Populate(reports[jumpNo:jumpNo+count], int64(count), req.Page, req.PageSize)

	return
}

func (s *TestResultService) GetLatest(siteId, productId uint) (summary serverDomain.TestReportSummary, err error) {
	workspaces, err := s.WorkspaceRepo.ListByProduct(siteId, productId)
	if err != nil {
		return
	}

	summaries := make([]serverDomain.TestReportSummary, 0, len(workspaces))
	for _, workspace := range workspaces {
		reportSeqs := analysisHelper.ListReport2(workspace.Path, 1)
		if len(reportSeqs) == 0 {
			continue
		}

		seq := reportSeqs[0]
		report, _, err1 := analysisHelper.ReadReportByWorkspaceSeq(workspace.Path, seq)
		if err1 != nil {
			continue
		}

		s := serverDomain.TestReportSummary{WorkspaceId: int(workspace.ID)}
		copier.Copy(&s, report)
		s.No = fmt.Sprintf("%d-%s", workspace.ID, seq)

		s.Seq = seq
		s.WorkspaceId = int(workspace.ID)
		s.WorkspaceName = workspace.Name
		if report.Total == 1 {
			_, s.TestScriptName = filepath.Split(report.FuncResult[0].Path)
		}

		summaries = append(summaries, s)
	}

	sort.Slice(summaries, func(i, j int) bool { return summaries[i].StartTime > summaries[j].StartTime })

	if len(summaries) > 0 {
		summary = summaries[0]
	}

	return
}

func (s *TestResultService) Get(workspaceId int, seq string) (report commDomain.ZtfReport, err error) {
	workspace, _ := s.WorkspaceRepo.Get(uint(workspaceId))
	report, _, err = analysisHelper.ReadReportByWorkspaceSeq(workspace.Path, seq)
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

	report, _, err := analysisHelper.ReadReportByWorkspaceSeq(workspace.Path, result.Seq)
	if err != nil {
		return
	}

	config := configHelper.LoadBySite(site)
	report.Name = result.Name
	err = zentaoHelper.CommitResult(report, result.ProductId, result.TaskId, config, nil)

	return
}

func (s *TestResultService) ZipLog(fileName string) (zipPath string, err error) {
	if fileName == "" {
		err = errors.New("file path is empty")
		return
	}
	zipPath = filepath.Join(commConsts.WorkDir, commConsts.DownloadServerPath, commConsts.ResultZip)
	path, _ := filepath.Split(fileName)
	fileUtils.RmDir(zipPath)
	fileUtils.ZipDir(zipPath, path)
	return
}

func (s *TestResultService) DownloadFromProxy(fileName string, workspaceId int, proxyPath string, pathMap map[string]string) (zipPath string, err error) {
	if fileName == "" {
		err = errors.New("file path is empty")
		return
	}
	if workspaceId == 0 {
		return
	}
	workspaceInfo, _ := s.WorkspaceService.Get(uint(workspaceId))
	if workspaceInfo.ID == 0 {
		err = errors.New("workspace not found")
		return
	}
	if workspaceInfo.ProxyId == 0 {
		return
	}
	url := ""
	if proxyPath != "" {
		url = proxyPath
	}
	if url == "" {
		err = errors.New("proxy path is empty")
		return
	}
	zipPath = filepath.Join(commConsts.WorkDir, commConsts.DownloadPath, commConsts.ResultZip)
	os.Remove(zipPath)
	fileUtils.Download(url+"api/v1/results/downloadLog?file="+fileName, zipPath)
	execLogDir := logUtils.GetLogDir(workspaceInfo.Path)
	fileUtils.Unzip(zipPath, execLogDir)
	paths, err := ioutil.ReadDir(execLogDir)
	if len(paths) == 0 {
		return
	}
	childrenDir := execLogDir + paths[0].Name()
	paths, err = ioutil.ReadDir(childrenDir)
	for _, path := range paths {
		fileUtils.CopyFile(fileUtils.AddSepIfNeeded(childrenDir)+path.Name(), execLogDir+path.Name())
		replaceProxyPath(execLogDir+path.Name(), fileName, execLogDir, pathMap)
	}
	fileUtils.RmDir(childrenDir)
	return execLogDir + commConsts.LogText, nil
}

func replaceProxyPath(fullPath, fileName, execLogDir string, pathMap map[string]string) (err error) {
	newContent := ""
	if fullPath[len(fullPath)-5:] == ".json" {
		report, err := analysisHelper.ReadReportByPath(fullPath)
		if err != nil {
			return err
		}
		for k, _ := range report.FuncResult {
			for proxyFilePath, localFilePath := range pathMap {
				report.FuncResult[k].Path = strings.ReplaceAll(report.FuncResult[k].Path, proxyFilePath, localFilePath)
				report.FuncResult[k].RelativePath = strings.ReplaceAll(report.FuncResult[k].RelativePath, proxyFilePath, localFilePath)
			}
		}
		newContentByte, err := json.Marshal(report)
		if err != nil {
			return err
		}
		newContent = string(newContentByte)
	} else {
		bytes, _ := ioutil.ReadFile(fullPath)
		newContent = strings.ReplaceAll(string(bytes), fileName[0:len(fileName)-10], execLogDir)
		for proxyFilePath, localFilePath := range pathMap {
			newContent = strings.ReplaceAll(newContent, proxyFilePath, localFilePath)
		}
	}
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	_, err = writer.Write([]byte(newContent))
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

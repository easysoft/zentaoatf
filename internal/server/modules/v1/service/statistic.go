package service

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	"github.com/jinzhu/copier"
)

type StatisticService struct {
	StatisticRepo    *repo.StatisticRepo `inject:""`
	WorkspaceService *WorkspaceService   `inject:""`
}

func NewStatisticService() *StatisticService {
	return &StatisticService{}
}

func (s *StatisticService) Get(id uint) (statistics model.Statistic, err error) {
	return s.StatisticRepo.Get(id)
}

func (s *StatisticService) GetByPath(path string) (statistics model.Statistic, err error) {
	return s.StatisticRepo.GetByPath(path)
}

func (s *StatisticService) Create(statistics model.Statistic) (id uint, isDuplicate bool, err error) {
	id, isDuplicate, err = s.StatisticRepo.Create(&statistics)

	return
}

func (s *StatisticService) Update(statistics model.Statistic) (isDuplicate bool, err error) {
	isDuplicate, err = s.StatisticRepo.Update(statistics)
	if isDuplicate || err != nil {
		return
	}
	return
}

func (s *StatisticService) Delete(id uint) error {
	return s.StatisticRepo.Delete(id)
}

func (s *StatisticService) UpdateStatistic(logPath string) (scriptPaths []string, err error) {
	pth := filepath.Join(logPath, commConsts.ResultJson)
	report, err := analysisHelper.ReadReportByPath(pth)
	scriptPaths = []string{}
	for _, statistic := range report.FuncResult {
		scriptPaths = append(scriptPaths, statistic.Path)
		exist, _ := s.GetByPath(statistic.Path)
		succ, fail := exist.Succ, exist.Fail
		if statistic.Status == "pass" {
			succ++
		} else {
			fail++
		}
		if exist.ID == 0 {
			_, _, err = s.StatisticRepo.Create(&model.Statistic{
				Path:     statistic.Path,
				Total:    succ + fail,
				Succ:     succ,
				Fail:     fail,
				FailLogs: pth,
			})
			return
		}
		err = s.StatisticRepo.UpdateStatistic(exist.Path, succ+fail, succ, fail, fmt.Sprintf("%s,%s", exist.FailLogs, pth))
		if err != nil {
			return
		}
	}
	return
}

func (s *StatisticService) GetFailureLogs(scriptPath string) (reports []serverDomain.TestReportSummary, err error) {
	statistic, err := s.GetByPath(scriptPath)
	if err != nil {
		return
	}
	logPath := strings.Split(statistic.FailLogs, ",")
	for _, path := range logPath {
		report, _ := analysisHelper.ReadReportByPath(path)
		workspace, _ := s.WorkspaceService.GetByPath(report.WorkspacePath)
		summary := serverDomain.TestReportSummary{
			WorkspaceId: int(workspace.ID),
		}
		copier.Copy(&summary, report)
		seq := strings.Replace(path, filepath.Join(workspace.Path, commConsts.LogDirName), "", -1)
		seq = strings.Replace(seq, commConsts.ResultJson, "", -1)
		seq = strings.Trim(seq, string(filepath.Separator))
		seq = analysisHelper.EncodeSeq(seq)

		summary.No = fmt.Sprintf("%d-%s", workspace.ID, seq)
		summary.Seq = seq
		summary.WorkspaceId = int(workspace.ID)
		summary.WorkspaceName = workspace.Name

		if report.Total == 1 && len(report.FuncResult) > 0 {
			_, summary.TestScriptName = filepath.Split(report.FuncResult[0].Path)
		}
		reports = append(reports, summary)
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].StartTime > reports[j].StartTime
	})
	return
}

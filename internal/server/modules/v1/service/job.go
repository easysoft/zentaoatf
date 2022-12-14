package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	channelUtils "github.com/easysoft/zentaoatf/pkg/lib/channel"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
)

var (
	channelMap sync.Map
)

type JobService struct {
	JobRepo *repo.JobRepo `inject:""`
}

func NewJobService() *JobService {
	return &JobService{}
}

func (s *JobService) Add(req serverDomain.ZentaoExecReq) (err error) {
	po := model.Job{
		Workspace: req.Workspace,
		Path:      req.Path,
		Ids:       req.Ids,

		Task:   req.Task,
		Retry:  1,
		Status: commConsts.JobCreated,
	}

	s.JobRepo.Save(&po)

	return
}

func (s *JobService) Start(po *model.Job) {
	ch := make(chan int, 1)
	channelMap.Store(po.ID, ch)

	req := s.genExecReq(*po)

	go func() {
		s.JobRepo.UpdateStatus(po, commConsts.JobInprogress, true, false)

		err := execHelper.Exec(nil, req, nil)

		s.JobRepo.UpdateStatus(po, commConsts.JobCompleted, false, true)

		// s.SubmitJobStatus(*po)

		s.SubmitExecResult(*po, err)

		if ch != nil {
			channelMap.Delete(po.ID)
			close(ch)
		}
	}()
}

func (s *JobService) Cancel(id uint) {
	taskInfo, _ := s.JobRepo.Get(id)

	if taskInfo.ID > 0 {
		s.JobRepo.SetCanceled(taskInfo)
	}

	s.stop(id)
}

func (s *JobService) Restart(po *model.Job) (ret bool) {
	s.stop(po.ID)
	s.Start(po)

	s.JobRepo.AddRetry(po)

	return
}

func (s *JobService) stop(id uint) {
	chVal, ok := channelMap.Load(id)

	if !ok || chVal == nil {
		return
	}

	channelMap.Delete(id)

	ch := chVal.(chan int)
	if ch != nil {
		if !channelUtils.IsChanClose(ch) {
			ch <- 1
		}

		ch = nil
	}
}

func (s *JobService) Check() (err error) {
	taskMap, _ := s.Query()

	toStartNewJob := false
	if len(taskMap.Inprogress) > 0 {
		runningJob := taskMap.Inprogress[0]

		if s.IsError(*runningJob) || s.IsTimeout(*runningJob) || s.isEmpty() {
			if s.NeedRetry(*runningJob) {
				s.Restart(runningJob)
			} else {
				s.JobRepo.UpdateStatus(runningJob, commConsts.JobFailed, false, true)
				s.SubmitJobStatus(*runningJob)

				toStartNewJob = true
			}
		}

	} else {
		toStartNewJob = true
	}

	if toStartNewJob && len(taskMap.Created) > 0 {
		newJob := taskMap.Created[0]

		s.Start(newJob)
	}

	return
}

func (s *JobService) List(status string) (jobs []model.Job, err error) {
	status = strings.TrimSpace(status)
	jobs, err = s.JobRepo.ListByStatus(status)

	return
}

func (s *JobService) Query() (ret serverDomain.JobQueryResp, err error) {
	//ret = serverDomain.JobQueryResp{
	//	Created:    make([]model.Job, 0),
	//	Inprogress: make([]model.Job, 0),
	//	Canceled:   make([]model.Job, 0),
	//	Completed:  make([]model.Job, 0),
	//	Failed:     make([]model.Job, 0),
	//}

	pos, _ := s.JobRepo.Query()

	for _, po := range pos {
		status := po.Status
		if status == commConsts.JobTimeout || status == commConsts.JobError {
			status = commConsts.JobInprogress
		}

		if status == commConsts.JobCreated {
			ret.Created = append(ret.Created, &po)
		} else if status == commConsts.JobInprogress {
			ret.Inprogress = append(ret.Inprogress, &po)
		} else if status == commConsts.JobCanceled {
			ret.Canceled = append(ret.Canceled, &po)
		} else if status == commConsts.JobCompleted {
			ret.Completed = append(ret.Completed, &po)
		} else if status == commConsts.JobFailed {
			ret.Failed = append(ret.Failed, &po)
		}
	}

	return
}

func (s *JobService) SubmitJobStatus(job model.Job) (err error) {
	status := serverDomain.ZentaoJobSubmitReq{
		Task:      job.Task,
		Status:    job.Status,
		StartTime: (*job.StartDate).Format("2006-01-02 15:04:05"),
		EndTime:   (*job.EndDate).Format("2006-01-02 15:04:05"),
		RetryTime: job.Retry,
		Error:     "",
		Data:      "",
	}

	config := commDomain.WorkspaceConf{
		Url: serverConfig.CONFIG.Server,
	}
	err = zentaoHelper.JobCommitResult(status, config)

	return
}

func (s *JobService) SubmitExecResult(job model.Job, execErr error) (err error) {
	result := serverDomain.ZentaoResultSubmitReq{
		Task: job.Task,
		Seq:  commConsts.ExecLogDir,
	}

	reportPth := filepath.Join(result.Seq, commConsts.ResultJson)
	var report commDomain.ZtfReport
	if fileUtils.FileExist(reportPth) {
		report, err = analysisHelper.ReadReportByPath(reportPth)
	} else {
		err = errors.New("case not found")
	}
	if err != nil && execErr == nil {
		execErr = err
	}

	config := commDomain.WorkspaceConf{
		Url: serverConfig.CONFIG.Server,
	}

	ret := serverDomain.ZentaoJobSubmitReq{
		Task:      job.Task,
		Status:    job.Status,
		StartTime: (*job.StartDate).Format("2006-01-02 15:04:05"),
		EndTime:   (*job.EndDate).Format("2006-01-02 15:04:05"),
		RetryTime: job.Retry,
		Error:     fmt.Sprintf("%v", execErr),
		Data:      report,
	}
	err = zentaoHelper.JobCommitResult(ret, config)

	return
}

func (s *JobService) genExecReq(po model.Job) (req serverDomain.ExecReq) {
	caseIds := make([]int, 0)
	for _, idStr := range strings.Split(po.Ids, ",") {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			caseIds = append(caseIds, id)
		}
	}

	dir := po.Path
	if !fileUtils.IsAbsolutePath(dir) {
		dir = filepath.Join(po.Workspace, dir)
	}

	caseIdMap := map[int]string{}
	scriptHelper.GetScriptByIdsInDir(dir, &caseIdMap)

	cases := scriptHelper.GetCaseByListInMap(caseIds, caseIdMap)

	commConsts.ExecFrom = commConsts.FromZentao
	req.Act = commConsts.ExecCase
	req.ScriptDirParamFromCmdLine = "."
	req.TestSets = append(req.TestSets, serverDomain.TestSet{
		WorkspacePath: po.Workspace,
		Cases:         cases,
	})

	return
}

func (s *JobService) IsError(po model.Job) bool {
	return po.Status == commConsts.JobError
}

func (s *JobService) IsTimeout(po model.Job) bool {
	dur := time.Now().Unix() - po.StartDate.Unix()
	// return dur > 3
	return po.Status == commConsts.JobInprogress && dur > commConsts.JobTimeoutTime
}

func (s *JobService) NeedRetry(po model.Job) bool {
	return po.Retry < commConsts.JobRetryTime
}

func (s *JobService) isEmpty() bool {
	length := 0

	channelMap.Range(func(key, value interface{}) bool {
		length++
		return true
	})

	return length == 0
}

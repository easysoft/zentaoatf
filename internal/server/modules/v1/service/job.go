package service

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	channelUtils "github.com/easysoft/zentaoatf/pkg/lib/channel"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
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

func (s *JobService) Start(po model.Job) {
	ch := make(chan int, 1)
	channelMap.Store(po.ID, ch)

	req := s.genExecReq(po)

	go func() {
		s.JobRepo.UpdateStatus(po.ID, commConsts.JobInprogress, true, false)

		execHelper.Exec(nil, req, nil)

		s.JobRepo.UpdateStatus(po.ID, commConsts.JobCompleted, false, true)

		s.SubmitResult(po)

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

	s.Stop(id)
}

func (s *JobService) Stop(id uint) {
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

func (s *JobService) Restart(po model.Job) (ret bool) {
	//s.Cancel(po.ID)
	s.Stop(po.ID)
	s.Start(po)

	s.JobRepo.AddRetry(po)

	return
}

func (s *JobService) Check() (err error) {
	taskMap, _ := s.Query()

	toStartNewJob := false
	if len(taskMap.Inprogress) > 0 {
		runningJob := taskMap.Inprogress[0]

		if s.IsError(runningJob) || s.IsTimeout(runningJob) || s.isEmpty() {
			if s.NeedRetry(runningJob) {
				s.Restart(runningJob)
			} else {
				s.JobRepo.SetFailed(runningJob)
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
			ret.Created = append(ret.Created, po)
		} else if status == commConsts.JobInprogress {
			ret.Inprogress = append(ret.Inprogress, po)
		} else if status == commConsts.JobCanceled {
			ret.Canceled = append(ret.Canceled, po)
		} else if status == commConsts.JobCompleted {
			ret.Completed = append(ret.Completed, po)
		} else if status == commConsts.JobFailed {
			ret.Failed = append(ret.Failed, po)
		}
	}

	return
}

func (s *JobService) SubmitResult(job model.Job) (err error) {
	result := serverDomain.ZentaoResultSubmitReq{
		Task: job.Task,
		Seq:  commConsts.ExecLogDir,
	}

	report, err := analysisHelper.ReadReportByPath(filepath.Join(result.Seq, commConsts.ResultJson))
	if err != nil {
		return
	}

	config := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)
	err = zentaoHelper.CommitResult(report, result.ProductId, result.TaskId, result.Task, config, nil)

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
	//return dur > 3
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

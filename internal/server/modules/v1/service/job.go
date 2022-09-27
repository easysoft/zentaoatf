package service

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/server/core/cache"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	"github.com/jinzhu/copier"
	"strconv"
	"time"
)

type JobService struct {
	JobRepo     *repo.JobRepo `inject:""`
	ExecService *ExecService  `inject:""`
}

func NewJobService() *JobService {
	return &JobService{}
}

func (s *JobService) Add(req serverDomain.JobReq) (err error) {
	job := model.Job{}
	copier.CopyWithOption(&job, req, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})

	s.JobRepo.Create(&job)

	return
}

func (s *JobService) Remove(req serverDomain.JobReq) (err error) {
	currJob := cache.GetCurrJob()

	if currJob.ID != 0 && currJob.ID == req.JobId {
		err = s.Stop()
	}

	s.JobRepo.Delete(req.JobId)

	return
}

func (s *JobService) Stop() (err error) {
	s.ExecService.Stop(nil)

	return
}

func (s *JobService) Check(currJob model.Job) (err error) {
	s.CheckJob()
	s.CheckTimeout()

	return
}

func (s *JobService) CheckJob() (err error) {
	job, err := s.JobRepo.QueryForExec()
	if err != nil {
		return
	}

	s.Run(job)

	return
}

func (s JobService) Run(job model.Job) (err error) {
	testSet := serverDomain.TestSet{
		Cases: s.convertIntToStrArr(job.CaseIds),
	}

	req := serverDomain.ExecReq{
		ProductId: job.ProductId,
		ModuleId:  job.ModuleId,
		SuiteId:   job.SuiteId,
		TaskId:    job.TaskId,
		TestSets:  []serverDomain.TestSet{testSet},
	}

	err = s.ExecService.Start(req, nil)
	if err != nil {
		return
	}

	s.JobRepo.UpdateProgressStatus(job.ID, commConsts.ProgressInProgress)

	return
}

func (s JobService) CheckTimeout() (err error) {
	currJob := cache.GetCurrJob()

	pos, err := s.JobRepo.ListByProgressStatus(commConsts.ProgressInProgress)
	if err != nil {
		return
	}

	for _, po := range pos {
		if time.Now().Unix()-po.StartTime.Unix() > commConsts.JobTimeout {
			if currJob.ID != 0 && currJob.ID == po.ID {
				err = s.Stop()
			}

			err = s.JobRepo.SetTimeout(currJob.ID)
		}
	}

	return
}

func (s *JobService) convertIntToStrArr(ids []int) (ret []string) {
	for _, item := range ids {
		ret = append(ret, strconv.Itoa(item))
	}

	return
}

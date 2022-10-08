package model

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"time"
)

type Job struct {
	BaseModel

	Name     string
	Priority int

	ScmAddress  string
	ScmAccount  string
	ScmPassword string
	ScmToken    string

	ProgressStatus commConsts.ProgressStatus
	BuildStatus    commConsts.BuildStatus

	StartTime     *time.Time
	TimeoutTime   *time.Time
	TerminateTime *time.Time
	ResultTime    *time.Time

	CaseIds   []int `json:"caseIds" gorm:"-"`
	ProductId int   `json:"productId"`
	ModuleId  int   `json:"moduleId"`
	SuiteId   int   `json:"suiteId"`
	TaskId    int   `json:"taskId"`
}

func NewJob() Job {
	task := Job{
		ProgressStatus: commConsts.ProgressCreated,
		BuildStatus:    commConsts.StatusCreated,
	}
	return task
}

func (Job) TableName() string {
	return "biz_job"
}

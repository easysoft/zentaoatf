package model

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"time"
)

type Job struct {
	BaseModel

	Name string `json:"name"`
	Desc string `json:"desc,omitempty"`

	Workspace string `json:"workspace"`
	Path      string `json:"path"`
	Ids       string `json:"ids"`

	Status commConsts.JobStatus `json:"status"`
	Retry  int                  `json:"retry"`

	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	TimeoutDate *time.Time `json:"timeoutDate,omitempty"`
	CancelDate  *time.Time `json:"cancelDate,omitempty"`

	Task int `json:"task"`
}

func (Job) TableName() string {
	return "biz_job"
}

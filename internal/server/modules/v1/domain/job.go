package serverDomain

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
)

type JobReq struct {
	JobId uint `json:"jobId"`

	Name      string `json:"name"`
	CaseIds   []uint `json:"caseIds"`
	ProductId int    `json:"productId"`
	ModuleId  int    `json:"moduleId"`
	SuiteId   int    `json:"suiteId"`
	TaskId    int    `json:"taskId"`
}

type JobResp struct {
	ZentaoId int                  `json:"zentaoId"`
	Status   commConsts.JobStatus `json:"status"` // Enums commConsts.JobStatus
}

type JobQueryResp struct {
	Created    []model.Job `json:"created"`
	Inprogress []model.Job `json:"inprogress"`

	Canceled  []model.Job `json:"canceled"`
	Completed []model.Job `json:"completed"`
	Failed    []model.Job `json:"failed"`
}

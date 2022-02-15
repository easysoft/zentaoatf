package command

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
)

type IndexModule struct {
	ProjectService    *service.ProjectService    `inject:""`
	SyncService       *service.SyncService       `inject:""`
	TestResultService *service.TestResultService `inject:""`
	TestBugService    *service.TestBugService    `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

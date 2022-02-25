package command

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
)

type IndexModule struct {
	SyncService *service.SyncService `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

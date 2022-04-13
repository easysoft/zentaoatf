package command

import (
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
)

type IndexModule struct {
	SyncService *service.SyncService `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

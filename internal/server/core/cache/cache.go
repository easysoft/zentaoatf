package cache

import (
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"sync"
)

const (
	IsRunning       = "IsRunning"
	LastLoopEndTime = "lastLoopEndTime"
	CurrJob         = "currJob"
)

var (
	SyncMap sync.Map
)

func GetCurrJob() (ret model.Job) {
	currJobObj, _ := SyncMap.Load(CurrJob)

	if currJobObj != nil {
		ret = currJobObj.(model.Job)
	}

	return
}

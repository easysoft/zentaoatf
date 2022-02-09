package scriptUtils

import (
	"sync"
)

const (
	KeyRunning = "isRunning"
)

var (
	SyncMap sync.Map
)

func SetRunning(val bool) {
	SyncMap.Store(KeyRunning, val)
}

func GetRunning() (val bool) {
	inf, ok := SyncMap.Load(KeyRunning)

	if ok {
		val = inf.(bool)
	}

	return
}

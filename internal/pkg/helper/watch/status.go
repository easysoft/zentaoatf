package watchHelper

import (
	"sync"
)

const (
	KeyWatching = "watching_"
)

var (
	SyncWatchMap sync.Map
)

func SetWatching(path string, val bool) {
	SyncWatchMap.Store(KeyWatching+path, val)
}

func GetWatching(path string) (val bool) {
	inf, ok := SyncWatchMap.Load(KeyWatching + path)

	if ok {
		val = inf.(bool)
	}

	return
}

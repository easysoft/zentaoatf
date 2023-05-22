package watchHelper

import (
	"encoding/json"
	"strings"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/pkg/consts"
	"github.com/fsnotify/fsnotify"
	"github.com/kataras/iris/v12/websocket"
)

func WatchFromReq(testSets []serverDomain.TestSet, wsMsg *websocket.Message) {
	paths := []string{}
	for _, testSet := range testSets {
		paths = append(paths, testSet.WorkspacePath)
	}
	Watch(paths, wsMsg)
}

func Watch(files []string, wsMsg *websocket.Message) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)

		watchEvent(watcher, wsMsg)
	}()

	for _, file := range files {
		if GetWatching(file) {
			continue
		}
		watcher.Add(file)
	}
	<-done
}

func watchEvent(watcher *fsnotify.Watcher, wsMsg *websocket.Message) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			//watch remove and create event
			if fsnotify.Remove == event.Op || fsnotify.Create == event.Op {
				if strings.Contains(event.Name, "log"+consts.FilePthSep) {
					return
				}

				//send websocket msg
				bytes, _ := json.Marshal(commDomain.WsResp{Category: "watch"})
				mqData := commDomain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
				websocketHelper.PubMsg(mqData)
			}
		case _, ok := <-watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

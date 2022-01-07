package cronUtils

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

var cronInst *cron.Cron

var taskFunc = make(map[string]func())

func GetCrontabInstance() *cron.Cron {
	if cronInst != nil {
		return cronInst
	}
	cronInst = cron.New()
	cronInst.Start()

	return cronInst
}

func AddTask(name string, schedule string, f func()) {
	if _, ok := taskFunc[name]; !ok {
		fmt.Printf("Add a new task %s", name)

		cInstance := GetCrontabInstance()
		cInstance.AddFunc(schedule, f)

		taskFunc[name] = f
	} else {
		fmt.Println("Don't add same task `" + name + "` repeatedly!")
	}
}

func Stop() {
	cronInst.Stop()
}

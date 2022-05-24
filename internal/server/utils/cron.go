package serverUtils

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

var cronInstace *cron.Cron

var taskFunc = make(map[string]func())

func GetCrontabInstance() *cron.Cron {
	if cronInstace != nil {
		return cronInstace
	}
	cronInstace = cron.New()
	cronInstace.Start()

	return cronInstace
}

func AddTaskFuc(name string, schedule string, f func()) {
	if _, ok := taskFunc[name]; !ok {
		fmt.Println("Add a new task:", name)

		cInstance := GetCrontabInstance()
		cInstance.AddFunc(schedule, f)

		taskFunc[name] = f
	} else {
		fmt.Println("Don't add same task `" + name + "` repeatedly!")
	}
}

func Stop() {
	cronInstace.Stop()
}

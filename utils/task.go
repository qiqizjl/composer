package utils

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var taskNumber = 0
var taskLock = sync.Mutex{}

func GetTask() int {
	return taskNumber
}

func ChangeTaskNumber(number int) {
	taskLock.Lock()
	logrus.Infoln("update task",taskNumber)
	taskNumber = taskNumber + number
	logrus.Infoln("update task",taskNumber)
	taskLock.Unlock()
}

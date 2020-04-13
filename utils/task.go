package utils

import "sync"

var taskNumber = 0
var taskLock = sync.Mutex{}

func GetTask() int {
	return taskNumber
}

func ChangeTaskNumber(number int) {
	taskLock.Lock()
	taskNumber = taskNumber + number
	taskLock.Unlock()
}

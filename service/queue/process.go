package queue

import (
	"strconv"
	"time"
)

func getProcessName(name string, num int) string {
	return name + ":" + strconv.Itoa(num)
}

func getDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}


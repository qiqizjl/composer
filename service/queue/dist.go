package queue

import (
	"composer/service/redis"
	"composer/service/upload"
	"composer/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"
)

func Dist(num int) {
	processName := getProcessName("dist", num)

	for {
		time.Sleep(1 * time.Second)
		job, err := redis.PopQueue(redis.Dist)
		if err != nil {
			// 无job
			continue
		}
		data := job[1]
		redis.RemoveQueue(redis.Dist, data, processName)
		jobData := new(utils.DistJob)
		err = json.Unmarshal([]byte(data), jobData)
		if err != nil {
			logrus.Errorln(processName, "job decode:", err.Error())
			continue
		}
		upload.Dist(processName, *jobData)

	}
}

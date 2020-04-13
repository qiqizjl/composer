package queue

import (
	"composer/service/redis"
	"composer/service/upload"
	"time"
)

func Provider(num int) {
	processName := getProcessName("provider", num)
	for {
		job, err := redis.PopQueue(redis.ProviderKey)
		if err != nil {
			// 无job
			time.Sleep(1 * time.Second)
			continue
		}
		path := job[1]
		redis.RemoveQueue(redis.ProviderKey, path, processName)
		upload.Provider(processName, path)

	}
}

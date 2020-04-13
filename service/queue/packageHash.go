package queue

import (
	"composer/service/redis"
	"composer/service/upload"
	"time"
)

func PackageHash(num int) {
	processName := getProcessName("packageHash", num)

	for {
		job, err := redis.PopQueue(redis.PackageHashFileKey)
		if err != nil {
			// 无job
			time.Sleep(1 * time.Second)
			continue
		}

		path := job[1]
		redis.RemoveQueue(redis.PackageHashFileKey, path, processName)
		// 调用上传代码
		upload.PackageHash(processName, path)
	}
}

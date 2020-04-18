package redis

import (
	"composer/utils"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	PackageHashFileKey = "mirror:packageHashFile"
	ProviderKey        = "mirror:providers"
	Dist               = "mirror:dist"
)

func inQueue(key, content, processName string) bool {
	key += ":queued"
	if hExists(key, content) {
		logrus.Debugln(processName, "Queued", key, content)
		return true
	}
	return false
}

func addQueue(key, content string) bool {
	key += ":queued"
	redisClient.HSet(key, content, time.Now().Unix())
	return true
}

func PopQueue(key string) ([]string, error) {
	timeout := 1 * time.Second
	key += ":queue"
	return redisClient.BRPop(timeout, key).Result()
}

func PushQueue(key, content, processName string) bool {
	if inQueue(key, content, processName) {
		return true
	}
	addQueue(key, content)
	key += ":queue"
	redisClient.LPush(key, content)
	return true
}

func RemoveQueue(key string, content string, processName string) {
	key += ":queued"
	logrus.Debugln(processName, "Queued Remove", key, content)
	redisClient.HDel(key, content)
}

func hExists(key, field string) bool {
	exist, err := redisClient.HExists(key, field).Result()
	if err != nil {
		logrus.Errorln(err.Error())
	}
	if exist {
		return true
	}
	return false
}

func zExists(key, field string) bool {
	_, err := redisClient.ZScore(key, field).Result()
	if err == nil {
		return true
	}
	if err != redis.Nil {
		logrus.Errorln(err.Error())
	}
	return false
}

func queueExists(key string) int64 {
	key += ":queue"
	num, err := redisClient.Exists(key).Result()
	if err != nil {
		num = 0
	}
	return num
}

func HasTask() bool {
	if queueExists(ProviderKey) != 0 {
		return true
	}
	if queueExists(PackageHashFileKey) != 0 {
		return true
	}
	if queueExists(Dist) != 0 {
		return true
	}
	if utils.GetTask() != 0 {
		return true
	}

	return false
}

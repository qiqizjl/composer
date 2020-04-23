package redis

import (
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

// IsSucceed 是否上传成功
func IsSucceed(key string, field string) bool {
	key += ":succeed"
	return hExists(key, field)
}

// makeSuccess 设置成功
func makeSuccess(key string, field string) {
	key += ":succeed"
	redisClient.HSet(key, field, time.Now().Unix())
}

// UploadSuccess 上传完成
func UploadSuccess(key string, path string) {
	// 设置完成
	makeSuccess(key, path)
	UpdateTime(key, path)
}

func UpdateTime(key string, field string) {
	redisClient.ZAdd(key, &redis.Z{Member: field, Score: float64(time.Now().Unix())})
}

func UploadSuccessTime(key string, path string, updateTime int) {
	// 设置完成
	makeSuccess(key, path)
	redisClient.ZAdd(key, &redis.Z{Member: path, Score: float64(updateTime)})
}

func GetFileList(key string, start, end int) []string {
	resp := make([]string, 0)
	result, err := redisClient.ZRangeByScore(key, &redis.ZRangeBy{
		Min: strconv.Itoa(start),
		Max: strconv.Itoa(end),
	}).Result()
	if err != nil {
		return resp
	}
	for _, data := range result {
		resp = append(resp, data)
	}
	return resp
}

func RemoveFile(key string, path string) {
	successKey := key + ":succeed"
	redisClient.HDel(successKey, path)
	fileKey := key
	redisClient.ZRem(fileKey, path)
}

func SetUpdateTime(value int64) {
	redisClient.Set(updateTimeKey, value, 0)
}
func GetUpdateTime() int {
	result, err := redisClient.Get(updateTimeKey).Int()
	if err != nil {
		return 0
	}
	return result
}

func AddRunTask(task string) {
	key := "mirrors:task"
	redisClient.HSet(key, task, time.Now().Unix())
}

func RemoveRunTask(task string) {
	key := "mirrors:task"
	redisClient.HDel(key, task)
}

func ClearRunTask() {
	key := "mirrors:task"
	redisClient.Del(key)
}

func AddDistSize(packageName string, size int) {
	key := Dist + ":size"
	redisClient.ZIncrBy(key, float64(size), packageName)
}

func ClearDistSize() {
	key := Dist + ":size"
	redisClient.Del(key)
}

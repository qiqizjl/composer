package _import

import (
	"composer/file"
	"composer/service/redis"
	"github.com/sirupsen/logrus"
)

func ImportDist(nextPage string) {
	logrus.Infoln("request", nextPage)
	result, err := file.DistFile.ListFile(nextPage, 100000)
	if err != nil {
		logrus.Errorln("request error", nextPage, err.Error())
		return
	}
	for item := range result {
		if !redis.IsSucceed(redis.Dist, item.Key) {
			redis.UploadSuccessTime(redis.Dist, item.Key, item.UpdateTime)
		}
		logrus.Infoln("update ", item.Key, item.UpdateTime)
	}

}

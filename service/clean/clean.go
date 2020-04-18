package clean

import (
	"composer/file"
	"composer/service/redis"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// 清理Metadata垃圾数据
func Metadata() {
	// 先从redis获取需要垃圾清理的列表
	day, _ := time.ParseDuration("-20m")

	cleanTime := int(time.Now().Add(day).Unix())
	_ = UpdateMetadataTime(false)
	providerList := redis.GetFileList(redis.ProviderKey, 0, cleanTime)
	if len(providerList) > 0 {
		for _, provider := range providerList {
			logrus.Infoln("remove provider", provider)
			cleanFile(provider)
			redis.RemoveFile(redis.ProviderKey, provider)
		}
	}

	packageHashList := redis.GetFileList(redis.PackageHashFileKey, 0, cleanTime)
	if len(packageHashList) > 0 {
		for _, packageHash := range packageHashList {
			logrus.Infoln("remove packageHash", packageHash)
			cleanFile(packageHash)
			redis.RemoveFile(redis.PackageHashFileKey, packageHash)
		}
	}
}

func Dist() {
	cleanTime := redis.GetUpdateTime()
	if cleanTime == 0 {
		return
	}
	distList := redis.GetFileList(redis.Dist, 0, cleanTime)
	if len(distList) > 0 {
		for _, dist := range distList {
			logrus.Infoln("remove dist", dist)
			file.DistFile.RemoveFile(dist)
			redis.RemoveFile(redis.Dist, dist)
		}
	}
}

func cleanFile(path string) {
	_ = os.Remove("./tmp/metadata/" + path)
	_ = file.MetaFile.RemoveFile(path)
}

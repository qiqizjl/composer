package clean

import (
	"composer/file"
	"composer/service/redis"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// 清理Metadata垃圾数据
func Metadata() {
	// 先从redis获取需要垃圾清理的列表
	day, _ := time.ParseDuration("-20m")

	cleanTime := int(time.Now().Add(day).Unix())

	fmt.Print(cleanTime)
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

func cleanFile(path string) {
	_ = os.Remove("./tmp/metadata/" + path)
	_ = file.MetaFile.RemoveFile(path)
}

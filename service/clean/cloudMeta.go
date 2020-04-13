package clean

import (
	"composer/file"
	"composer/service/redis"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func CloudMeta() {
	result, err := file.MetaFile.ListFile("", 100000)
	if err != nil {
		logrus.Errorln("request error", err.Error())
		return
	}
	for item := range result {
		if item.Key == "packages.json" {
			continue
		}
		key := redis.PackageHashFileKey
		// 是否是
		if strings.Contains(item.Key, "p/provider") {
			key = redis.ProviderKey
		}
		if !redis.IsSucceed(key, item.Key) {
			file.MetaFile.RemoveFile(item.Key)
			fmt.Println("clean remote", item.Key)
		}

	}

}

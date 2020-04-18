package clean

import (
	"composer/file"
	"composer/service/redis"
	"composer/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

// 清理Metadata垃圾数据
func Metadata() {
	// 先从redis获取需要垃圾清理的列表
	day, _ := time.ParseDuration("-20m")

	cleanTime := int(time.Now().Add(day).Unix())
	_ = UpdateMetadataTime()
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

func UpdateMetadataTime() error {
	packages, err := utils.GetPackages()
	fmt.Println(err)
	if err != nil {
		return err
	}
	providerHashList := packages.ProvidersList()
	wait := sync.WaitGroup{}
	for _, path := range providerHashList {
		wait.Add(1)
		path := path
		go func() {
			logrus.Infoln("update provider hash", path)
			redis.UpdateTime(redis.ProviderKey, path)
			updateProvider(path)
			wait.Done()
		}()
	}
	wait.Wait()
	return nil
}

func updateProvider(path string) {
	provider, err := utils.GetProviderInfo(path)
	if err != nil {
		return
	}
	for path := range provider.PackageList() {
		logrus.Infoln("update package", path)
		//流处理
		redis.UpdateTime(redis.PackageHashFileKey, path)
		//updatePackage(path)
	}
}

func updatePackage(path string) {
	packageInfo, err := utils.GetPackageInfo(path)
	if err != nil {
		return
	}
	for distJob := range packageInfo.GetDistPath() {
		redis.UpdateTime(redis.Dist, distJob.Path)
	}
}

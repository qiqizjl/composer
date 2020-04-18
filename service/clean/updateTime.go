package clean

import (
	"composer/service/redis"
	"composer/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

func UpdateMetadataTime(updateDist bool) error {
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
			updateProvider(path, updateDist)
			wait.Done()
		}()
	}
	wait.Wait()
	return nil
}

func updateProvider(path string, updateDist bool) {
	provider, err := utils.GetProviderInfo(path)
	if err != nil {
		logrus.Infoln("get provider error", path)
		return
	}
	for path := range provider.PackageList() {
		logrus.Infoln("update package", path)
		//流处理
		redis.UpdateTime(redis.PackageHashFileKey, path)
		if updateDist {
			updatePackage(path)
		}
	}
}

func updatePackage(path string) {
	packageInfo, err := utils.GetPackageInfo(path)
	if err != nil {
		logrus.Errorln("get package dist error", path, err.Error())
		return
	}
	for distJob := range packageInfo.GetDistPath() {
		redis.UpdateTime(redis.Dist, distJob.Path)
	}
}

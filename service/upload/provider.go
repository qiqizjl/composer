package upload

import (
	"composer/file"
	"composer/service/http"
	"composer/service/redis"
	"composer/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func Provider(processName string, path string) {
	utils.ChangeTaskNumber(1)
	defer utils.ChangeTaskNumber(-1)
	if redis.IsSucceed(redis.ProviderKey, path) {
		logrus.Println(processName, "file local exist:", path)
		redis.UpdateTime(redis.ProviderKey, path)

		return
	}
	//if file.MetaFile.IsFile(path) {
	//	//文件存在
	//	logrus.Println(processName, "file cloud exist:", path)
	//	redis.UploadSuccess(redis.ProviderKey, path)
	//	return
	//}

	resp, err := http.PackagistGet(path, processName)
	if err != nil {
		//服务异常 下一个
		logrus.Errorln(processName, path, " err:", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// 不等于200
		logrus.Errorln(processName, path, " status:", resp.StatusCode)
		return
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln(processName, path, " get:", resp.StatusCode)
		return
	}
	content, _ = utils.Decode(content)

	_, err = file.MetaFile.UploadFile(path, content)
	if err != nil {
		logrus.Errorln(processName, path, " upload error:", err.Error())
		return
	}

	err = utils.StoreMetadata(path, content)
	if err != nil {
		logrus.Infoln("store ", path, " error ", err.Error())
	}

	redis.UploadSuccess(redis.ProviderKey, path)

	packageList := make(map[string]interface{})
	err = json.Unmarshal(content, &packageList)
	if err != nil {
		logrus.Errorln(processName, path, "json_decode err: ", err.Error())
		return
	}
	dispatchPackages(packageList["providers"], processName)
}

func dispatchPackages(packageList interface{}, processName string) {
	for packageName, value := range packageList.(map[string]interface{}) {
		for _, hash := range value.(map[string]interface{}) {
			path := "p/" + packageName + "$" + hash.(string) + ".json"

			if redis.IsSucceed(redis.PackageHashFileKey, path) {
				redis.UpdateTime(redis.PackageHashFileKey, path)
				logrus.Traceln(processName, "file local exist:", path)
				continue
			}
			redis.PushQueue(redis.PackageHashFileKey, path, processName)
		}

	}

}

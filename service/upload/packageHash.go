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

// PackageHash 上传
func PackageHash(processName string, path string) {
	//先判断文件是否存在
	if redis.IsSucceed(redis.PackageHashFileKey, path) {
		logrus.Println(processName, "file local exist:", path)
		redis.UpdateTime(redis.PackageHashFileKey, path)
		return
	}
	//if file.MetaFile.IsFile(path) {
	//	//文件存在
	//	logrus.Println(processName, "file cloud exist:", path)
	//	redis.UploadSuccess(redis.PackageHashFileKey, path)
	//	return
	//}
	utils.ChangeTaskNumber(1)
	defer utils.ChangeTaskNumber(-1)

	//远程获取文件
	resp, err := http.PackagistGet(path, processName)
	// 下载失败
	if err != nil {
		logrus.Errorln(processName, path, " download error:", err.Error())
		return
	}
	// 关闭连接
	defer resp.Body.Close()
	// 非200
	if resp.StatusCode != 200 {
		logrus.Errorln(processName, path, " download error:", resp.Status)
		return
	}

	fileData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln(processName, path, " io read error:", err.Error())
		return
	}
	fileData, _ = utils.Decode(fileData)

	// 先分发
	packageList := make(map[string]interface{})
	err = json.Unmarshal(fileData, &packageList)
	if err != nil {
		logrus.Errorln(processName, path, " json decode error:", err.Error())
		return
	}
	dispatchDist(packageList["packages"], processName, path)

	// 后上传
	_, err = file.MetaFile.UploadFile(path, fileData)
	if err != nil {
		logrus.Errorln(processName, path, "upload error:", err.Error())
		return
	}
	err = utils.StoreMetadata(path, fileData)
	if err != nil {
		logrus.Infoln("store ", path, " error ", err.Error())
	}
	// 上传成功就写入redis
	redis.UploadSuccess(redis.PackageHashFileKey, path)

}

// dispatchDist 调用Dist文件
func dispatchDist(packages interface{}, processName string, path string) {
	list, ok := packages.(map[string]interface{})
	if !ok {
		return
	}
	for packageName, value := range list {
		for version, versionContent := range value.(map[string]interface{}) {
			dist, ok := versionContent.(map[string]interface{})["dist"].(map[string]interface{})
			if !ok {
				continue
			}
			if dist["reference"] == nil {
				logrus.Errorln(processName, path, version, "dist not found :", path, dist)
				continue
			}
			// 同步版本
			versionPath := packageName + "/" + dist["reference"].(string) + "." + dist["type"].(string)
			// 先判断文件是否存在
			if redis.IsSucceed(redis.Dist, versionPath) {
				// 文件存在 刷新文件时间
				redis.UpdateTime(redis.Dist, versionPath)
				continue
			}

			distJob := utils.DistJob{}
			distJob.Path = versionPath
			distJob.ContentURL = dist["url"].(string)
			jobContent, _ := json.Marshal(distJob)
			redis.PushQueue(redis.Dist, string(jobContent), processName)
		}
	}
}

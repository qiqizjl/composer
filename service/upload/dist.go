package upload

import (
	"composer/file"
	"composer/service/http"
	"composer/service/redis"
	"composer/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

//Dist 上传静态资源
func Dist(processName string, jobData utils.DistJob) {
	utils.ChangeTaskNumber(1)
	defer utils.ChangeTaskNumber(-1)
	nowRunTaskKey := fmt.Sprintf("dist_%s",jobData.Path)
	redis.AddRunTask(nowRunTaskKey)
	defer redis.RemoveRunTask(nowRunTaskKey)

	//先判断文件是否存在
	if redis.IsSucceed(redis.Dist, jobData.Path) {
		logrus.Println(processName, "file local exist:", jobData.Path)
		redis.UpdateTime(redis.Dist, jobData.Path)
		return
	}
	if file.DistFile.IsFile(jobData.Path) {
		//文件存在
		logrus.Println(processName, "file cloud exist:", jobData.Path)
		redis.UploadSuccess(redis.Dist, jobData.Path)
		return
	}

	// push文件

	resp, err := http.DistGet(jobData.ContentURL, processName)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		logrus.Infoln(processName, jobData.Path, "download error:", err.Error())
		return
	}

	if resp.StatusCode != 200 {
		logrus.Infoln(processName, jobData.Path, "download error:", resp.StatusCode, " message", resp.Status)
		return
	}
	_, err = file.DistFile.UploadFileIO(jobData.Path, resp.Body, -1)
	if err != nil {
		logrus.Errorln(processName, jobData, "upload error:", err.Error())
		return
	}
	logrus.Errorln(processName, "upload success:", jobData)
	redis.UploadSuccess(redis.Dist, jobData.Path)
	// 人工GC
	runtime.GC()
}

package queue

import (
	"bytes"
	"composer/file"
	"composer/service/http"
	"composer/service/redis"
	"composer/utils"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
	"time"
)

var packagesJsonCache []byte
var packagistLastModified = ""

func Package(num int) {
	processName := getProcessName("package", num)

	for true {
		time.Sleep(1 * time.Second)
		if !redis.HasTask() {
			break
		}
		logrus.Infoln(processName, "wait task end")
	}

	for {
		time.Sleep(1 * time.Second)
		resp, err := http.PackagistGet("packages.json", processName)
		if err != nil {
			logrus.Errorln(processName, "download error packages.json", err.Error())
			continue
		}
		// Status code must be 200
		if resp.StatusCode != 200 {
			logrus.Errorln(processName, "packages.json", resp.Status)
			continue
		}

		// Read data stream from body
		content, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			logrus.Errorln(processName, "packages.json", err.Error())
			continue
		}


		content, _ = utils.Decode(content)

		// Get Last-Modified field
		if bytes.Equal(packagesJsonCache, content) {
			fmt.Println(processName, "Update to date: packages.json")
			continue
		}
		getSourceTime := time.Now()
		packagesJsonCache = content

		packagistLastModified = resp.Header["Last-Modified"][0]

		var packagesJson = make(map[string]interface{})

		// JSON Decode
		err = json.Unmarshal(content, &packagesJson)
		if err != nil {
			logrus.Errorln("Error: %s\n", err.Error())
			continue
		}
		downloadProviders(packagesJson["provider-includes"], processName)
		for true {
			time.Sleep(1 * time.Second)
			if !redis.HasTask() {
				break
			}
			logrus.Infoln(processName, "Synchronization task is not completed, check again in 1 second.")
		}
		packagistLastUpdateTime, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", packagistLastModified)
		updateTime := map[string]interface{}{}

		updateTime["mirrors-last-update"] = time.Now().Format("2006-01-02 15:04:05")
		updateTime["source-last-update"] = packagistLastUpdateTime.In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")
		updateTime["mirrors-get-update"] = getSourceTime.Format("2006-01-02 15:04:05")

		//getSourceTime
		packagesJson["update_at"] = updateTime

		packagesJson["mirrors"] = []map[string]interface{}{
			{
				"dist-url":  viper.GetString("file.download_domain.dist") + "%package%/%reference%.%type%",
				"preferred": true,
			},
		}

		// Json Encode
		content, _ = json.Marshal(packagesJson)
		_, err = file.MetaFile.UploadFile("packages.json", content)
		if err != nil {
			logrus.Infoln("update packagejson error ", err.Error())

		}
		logrus.Infoln("update packagejson", string(content))

		logrus.Errorln("update composer mirrors success")
	}
}

func downloadProviders(providerList interface{}, processName string) {
	for provider, value := range providerList.(map[string]interface{}) {
		for _, hash := range value.(map[string]interface{}) {

			path := strings.Replace(provider, "%hash%", hash.(string), -1)

			if redis.IsSucceed(redis.ProviderKey, path) {
				redis.UpdateTime(redis.ProviderKey, path)
				logrus.Traceln(processName, "file local exist:", path)
				continue
			}
			redis.PushQueue(redis.ProviderKey, path, processName)
		}

	}
}

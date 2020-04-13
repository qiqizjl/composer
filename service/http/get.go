package http

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func packagistUrl(url string) string {
	return viper.GetString("system.source") + url
}

func PackagistGet(url string, processName string) (*http.Response, error) {
	logrus.Infoln(processName,"Get Downloading", url)
	return GetRemoteGzip(packagistUrl(url))
}

func DistGet(url string, processName string) (*http.Response, error) {
	logrus.Infoln(processName, "Get Dist", url)
	return GetRemoteData(url)
}

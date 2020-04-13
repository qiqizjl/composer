package http

import (
	"composer/version"
	"crypto/tls"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"
)

var client *http.Client

func init() {
}

func getClient() *http.Client {
	if client == nil {
		// 初始化HTTP基础类
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		// HTTP代理
		if viper.IsSet("system.proxy") {
			proxy, _ := url.Parse(viper.GetString("system.proxy"))
			transport.Proxy = http.ProxyURL(proxy)
		}
		// HTTP客户端
		client = &http.Client{
			Transport: transport,
			Timeout:   time.Second * viper.GetDuration("system.timeout"),
		}
	}
	return client
}

func GetRemoteData(url string) (*http.Response, error) {
	request, _ := http.NewRequest("GET", url, nil)
	hostname, _ := os.Hostname()
	request.Header.Set("User-Agent",
		fmt.Sprintf(
			"Composer-Mirrors-Spider/%s(%s)  %s(%s)/%s",
			version.VERSION,
			version.BUILDDATE,
			runtime.Version(),
			runtime.GOOS,
			hostname,
		),
	)
	resp, err := getClient().Do(request)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetRemoteGzip(url string)  (*http.Response, error)  {
	request, _ := http.NewRequest("GET", url, nil)
	hostname, _ := os.Hostname()
	request.Header.Set("User-Agent",
		fmt.Sprintf(
			"Composer-Mirrors-Spider/%s(%s)  %s(%s)/%s",
			version.VERSION,
			version.BUILDDATE,
			runtime.Version(),
			runtime.GOOS,
			hostname,
		),
	)
	request.Header.Add("Content-Encoding", "gzip")
	request.Header.Add("Accept-Encoding", "gzip")
	resp, err := getClient().Do(request)
	if err != nil {
		return nil, err
	}
	return resp, err
}

package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

type DistJob struct {
	Path       string `json:"path"`
	ContentURL string `json:"content_url"`
}

//PackagistPackages packages.json 解析
type PackagistPackages struct {
	MetadataURL      string                  `json:"metadata-url"`
	Notify           string                  `json:"notify"`
	NotifyBatch      string                  `json:"notify-batch"`
	ProvidersApi     string                  `json:"providers-api"`
	ProvidersURL     string                  `json:"providers-url"`
	Search           string                  `json:"search"`
	Packages         []interface{}           `json:"packages"`
	ProviderIncludes map[string]ProviderHash `json:"provider-includes"`
}

func (packages *PackagistPackages) ProvidersList() []string {
	hashList := make([]string, 0)
	for fileName, hash := range packages.ProviderIncludes {
		fileName = strings.Replace(fileName, "%hash%", hash.Sha256, -1)
		hashList = append(hashList, fileName)
	}
	return hashList
}

// Provide provide文件解析
type PackagistProvide struct {
	Providers map[string]ProviderHash `json:"providers"`
}

// PackagistPackage 包信息
type PackagistPackage struct {
	Packages map[string]map[string]PackagistPackageVersion `json:"packages"`
}

func (packagistPackage *PackagistPackage) GetDistPath() chan DistJob {
	resList := make(chan DistJob)
	go func() {
		reg := regexp.MustCompile(viper.GetString("system.ignore"))
		for packageName, packageInfo := range packagistPackage.Packages {
			if reg.MatchString(packageName) {
				logrus.Info(packageName," ignore dist update")
				continue
			}
			for _, versionInfo := range packageInfo {
				if versionInfo.Dist.Reference == "" {
					// 为空不循环
					continue
				}
				versionPath := packageName + "/" + versionInfo.Dist.Reference + "." + versionInfo.Dist.Type
				distJob := DistJob{}
				distJob.Path = versionPath
				distJob.ContentURL = versionInfo.Dist.URL
				resList <- distJob
			}
		}
		close(resList)
	}()
	return resList
}

func (packagistProvide *PackagistProvide) PackageList() chan string {
	resList := make(chan string)
	go func() {
		reg := regexp.MustCompile(viper.GetString("system.ignore"))
		for packageName, m := range packagistProvide.Providers {
			if reg.MatchString(packageName) {
				logrus.Info(packageName," ignore update")
				continue
			}
			path := "p/" + packageName + "$" + m.Sha256 + ".json"
			select {
			case resList <- path:
			}
		}
		close(resList)
	}()
	return resList
}

//PackagistPackageVersion 版本信息
type PackagistPackageVersion struct {
	Dist struct {
		Type      string `json:"type"`
		URL       string `json:"url"`
		Reference string `json:"reference"`
		Shasum    string `json:"shasum"`
	}
}

type ProviderHash struct {
	Sha256 string `json:"sha256"`
}

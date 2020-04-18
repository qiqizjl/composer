package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// StoreMetadata 存储文件
func StoreMetadata(path string, data []byte) error {
	localPath := getMetadataPath(path)
	paths, _ := filepath.Split(localPath)

	err := os.MkdirAll(paths, 0755)
	if err != nil {
		return err
	}
	//写入本地文件
	err = ioutil.WriteFile(localPath, data, 0644)
	return err
}

// getMetadataPath 获得文件路径
func getMetadataPath(path string) string {
	return "./tmp/metadata/" + path
}

func GetMetadata(path string) ([]byte, error) {
	localPath := getMetadataPath(path)
	return ioutil.ReadFile(localPath)
}

//GetPackages 获得本地packages.json信息
func GetPackages() (*PackagistPackages, error) {
	content, err := GetMetadata("packages.json")
	if err != nil {
		return nil, err
	}
	result := &PackagistPackages{}
	err = json.Unmarshal(content, result)
	fmt.Print(result)

	return result, err
}

func GetProviderInfo(path string) (*PackagistProvide, error) {
	content, err := GetMetadata(path)
	if err != nil {
		return nil, err
	}
	result := &PackagistProvide{}
	err = json.Unmarshal(content, result)
	return result, err
}

func GetPackageInfo(path string)(*PackagistPackage, error)  {
	content, err := GetMetadata(path)
	if err != nil {
		return nil, err
	}
	result := &PackagistPackage{}
	err = json.Unmarshal(content, result)
	return result, err
}
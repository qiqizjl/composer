package file

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
)

type Adapter interface {
	// 文件是否存在
	IsFile(file string) bool
	// 上传文件
	UploadFile(file string, fileData []byte) (bool, error)
	// 通过IO接口上传文件
	UploadFileIO(file string, fileData io.Reader, fileLen int64) (bool, error)
	// 刷新CDN缓存
	RefreshFile(file string) (bool, error)
	// 删除文件
	RemoveFile(file string) bool
	// 列出文件
	ListFile(page string, pageSize int) (chan *ListItem, error)
}

type ListItem struct {
	Key        string
	Size       int
	UpdateTime int
	Remark     string
}

type List struct {
	Page     string
	NextPage string
	ListItem []ListItem
}

var (
	MetaFile Adapter
	DistFile Adapter
)

func InitFile() {
	MetaFile = getFileAdapter("meta")
	DistFile = getFileAdapter("dist")
}

func getFileAdapter(bucket string) Adapter {
	bucketName := viper.GetString(fmt.Sprintf("file.bucket.%s", bucket))
	switch viper.GetString("file.adapter") {
	case "qiniu":
		return newQiniu(bucketName)
	}

	return nil
}

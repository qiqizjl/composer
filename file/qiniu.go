package file

import (
	"bytes"
	"context"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/viper"
	"io"
)

type qiniuResult struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

type qiniu struct {
	mac    *qbox.Mac
	bucket string
}

func newQiniu(bucket string) *qiniu {
	qiniu := &qiniu{
		mac: qbox.NewMac(
			viper.GetString("file.config.access_key"),
			viper.GetString("file.config.secert_key"),
		),
		bucket: bucket,
	}
	return qiniu
}

func (q *qiniu) IsFile(file string) bool {
	bucketManager := storage.NewBucketManager(q.mac, q.getConfig())
	_, err := bucketManager.Stat(q.bucket, file)
	if err != nil {
		return false
	}
	return true
}

func (q *qiniu) UploadFile(file string, fileData []byte) (bool, error) {
	formUploader := storage.NewFormUploader(q.getConfig())
	result := qiniuResult{}
	err := formUploader.Put(context.Background(), &result, q.getUploadToken(file), file, bytes.NewBuffer(fileData), int64(len(fileData)), nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (q *qiniu) UploadFileIO(file string, fileData io.Reader, fileLen int64) (bool, error) {
	resumeUploader := storage.NewResumeUploader(q.getConfig())
	result := qiniuResult{}
	err := resumeUploader.PutWithoutSize(context.Background(), &result, q.getUploadToken(file), file, fileData, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (q *qiniu) getUploadToken(file string) string {
	putPolicy := storage.PutPolicy{
		Scope: q.bucket + ":" + file,
	}
	return putPolicy.UploadToken(q.mac)
}

func (q *qiniu) getConfig() *storage.Config {
	return &storage.Config{
		UseCdnDomains: true,
	}
}

func (q *qiniu) RefreshFile(file string) (bool, error) {
	return false, nil
}

func (q *qiniu) RemoveFile(file string) bool {
	bucketManager := storage.NewBucketManager(q.mac, q.getConfig())
	err := bucketManager.Delete(q.bucket, file)
	if err != nil {
		return false
	}
	return true
}

func (q *qiniu) ListFile(page string, pageSize int) (chan *ListItem, error) {
	bucketManager := storage.NewBucketManager(q.mac, q.getConfig())
	entries, err := bucketManager.ListBucket(q.bucket, "", "", page)
	if err != nil {
		return nil, err
	}
	retCh := make(chan *ListItem)
	go func() {
		for listItem := range entries {
			item := &ListItem{
				Key:        listItem.Item.Key,
				Size:       int(listItem.Item.Fsize),
				UpdateTime: int(listItem.Item.PutTime / 10000000),
				Remark:     listItem.Marker,
			}
			retCh <- item
		}
		close(retCh)
	}()

	return retCh, nil
}

//if err != nil {
//	return nil, err
//}
//marker := ""
//result := &List{
//	Page: page,
//}
//i := 0
//for listItem := range entries {
//	i++
//	result.ListItem = append(result.ListItem, ListItem{
//		Key:        listItem.Item.Key,
//		Size:       int(listItem.Item.Fsize),
//		UpdateTime: int(listItem.Item.PutTime / 10000000),
//	})
//	marker = listItem.Marker
//	if i >= pageSize {
//		break
//	}
//}
//result.NextPage = marker
//
//return result, nil

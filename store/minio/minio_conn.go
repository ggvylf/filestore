package store

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var mc *minio.Client

func init() {
	mc = NewMinioClient()

}

func NewMinioClient() *minio.Client {
	endpoint := "127.0.0.1:9000"
	aceessKey := "KwoLR7sWIdp8LZAt"
	secretAccessKey := "6sZAYToCkY7Uhl0hZbJjIgewReVEyLXt"
	userSSL := false
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(aceessKey, secretAccessKey, ""),
		Secure: userSSL,
	})
	if err != nil {
		fmt.Println("conn minio failed,err=", err)
	}

	return mc

}

func GetMC() *minio.Client {
	return mc
}

// 从oss获取对象下载地址
func DownloadUrl(filehash, filename string) string {

	// 获取mc实例

	ctx := context.Background()
	mc := GetMC()
	bucket := "userfile"
	ossName := "/minio" + "/" + filehash
	// path := "/userfile" + ossName

	// 从oss中获取object的下载地址
	// 添加响应header，注意响应header的名字跟标准header不同，以response-开头
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(filename)))

	url, err := mc.PresignedGetObject(ctx, bucket, ossName, time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println("获取oss下载地址失败")
		return ""
	}
	return url.String()

}

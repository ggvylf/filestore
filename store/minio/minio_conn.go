package store

import (
	"fmt"

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

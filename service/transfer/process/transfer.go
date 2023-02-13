package process

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/ggvylf/filestore/mq"
	dbcli "github.com/ggvylf/filestore/service/dbproxy/client"
	store "github.com/ggvylf/filestore/store/minio"
	"github.com/minio/minio-go/v7"
)

// 消费mq回调函数
func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))

	// 解析msg数据
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 打开本地临时文件
	file, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer file.Close()

	// 写入oss
	ctx := context.Background()

	_, err = store.GetMC().PutObject(
		ctx,
		pubData.Bucket,
		pubData.DestLocation,
		bufio.NewReader(file),
		pubData.FileSize,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 更新文件地址
	resp, err := dbcli.UpdateFileLocation(pubData.FileHash, pubData.DestLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !resp.Suc {
		log.Println("更新数据库异常，请检查:" + pubData.FileHash)
		return false
	}

	return true
}

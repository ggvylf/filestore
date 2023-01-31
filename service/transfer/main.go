package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/ggvylf/filestore/config"
	dblayer "github.com/ggvylf/filestore/db"
	"github.com/ggvylf/filestore/mq"
	store "github.com/ggvylf/filestore/store/minio"
	"github.com/minio/minio-go/v7"
)

// ProcessTransfer : 处理文件转移
func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	ctx := context.Background()

	_, err = store.GetMC().PutObject(
		ctx,
		pubData.Bucket,
		pubData.DestLocation,
		bufio.NewReader(fin),
		pubData.FileSize,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	_ = dblayer.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	return true
}

func main() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer)
}

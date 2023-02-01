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
	_ = dblayer.UpdateFmAddr(pubData.FileHash, pubData.DestLocation)

	return true
}

func main() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")

	// 消费mq
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer)
}

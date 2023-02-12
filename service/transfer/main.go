package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ggvylf/filestore/config"
	dblayer "github.com/ggvylf/filestore/db"
	"github.com/ggvylf/filestore/mq"
	store "github.com/ggvylf/filestore/store/minio"
	"github.com/minio/minio-go/v7"
	"go-micro.dev/v4"
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

// 文件传输服务
func startTransferService() {
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

// 微服务框架
func StartRpcService() {
	service := micro.NewService(
		micro.Name("go.micro.service.transfer"), // 服务名称
		micro.RegisterTTL(time.Second*10),       // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5),   // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
	)
	service.Init()

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func main() {

	// 并发读取mq
	go startTransferService()

	// 微服务相关
	StartRpcService()

}

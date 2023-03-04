package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ggvylf/filestore/config"
	"github.com/ggvylf/filestore/mq"
	dbproxy "github.com/ggvylf/filestore/service/dbproxy/client"
	"github.com/ggvylf/filestore/service/transfer/process"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

// 文件传输服务
func startTransferService() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")

	// 初始化mq连接
	mq.Init()

	// 消费mq
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		process.ProcessTransfer)

}

// 微服务框架
// transfer从mq消费，没有rpc调用
func StartRpcService() {

	consul := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Server(grpc.NewServer()),
		micro.Name("go.micro.service.transfer"), // 服务名称
		micro.RegisterTTL(time.Second*10),       // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5),   // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		micro.Registry(consul),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init()

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

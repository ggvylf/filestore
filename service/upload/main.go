package main

import (
	"fmt"
	"time"

	"github.com/ggvylf/filestore/mq"
	dbproxy "github.com/ggvylf/filestore/service/dbproxy/client"
	cfg "github.com/ggvylf/filestore/service/upload/config"
	handler "github.com/ggvylf/filestore/service/upload/handler"
	upProto "github.com/ggvylf/filestore/service/upload/proto"
	"github.com/ggvylf/filestore/service/upload/route"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

func startRpcService() {

	consul := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Server(grpc.NewServer()),
		micro.Name("go.micro.service.upload"), // 服务名称
		micro.RegisterTTL(time.Second*10),     // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		micro.Registry(consul),
	)
	// 初始化服务
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init()

	// 初始化mq client
	mq.Init()

	// 服务注册
	upProto.RegisterUploadServiceHandler(service.Server(), new(handler.Upload))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startApiService() {
	router := route.Router()
	router.Run(cfg.UploadServiceHost)
}

func main() {
	// api 服务
	go startApiService()

	// rpc 服务
	startRpcService()
}

package main

import (
	"fmt"
	"time"

	cfg "github.com/ggvylf/filestore/service/download/config"
	"github.com/ggvylf/filestore/service/download/route"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"

	dbproxy "github.com/ggvylf/filestore/service/dbproxy/client"
	dlRpc "github.com/ggvylf/filestore/service/download/handler"
	dlProto "github.com/ggvylf/filestore/service/download/proto"
)

func startRpcService() {

	consul := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Server(grpc.NewServer()),
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(consul),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init()

	// 注册服务到注册中心
	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// web接口
func startApiService() {
	router := route.Router()
	router.Run(cfg.DownloadServiceHost)
}

func main() {
	// api 服务
	go startApiService()

	// rpc 服务
	startRpcService()
}

package main

import (
	"log"
	"time"

	"github.com/ggvylf/filestore/service/account/handler"
	proto "github.com/ggvylf/filestore/service/account/proto"
	dbproxy "github.com/ggvylf/filestore/service/dbproxy/client"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

var (
	service_name = "go.micro.service.user"
)

func main() {
	consul := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	// 创建一个微服务
	service := micro.NewService(
		micro.Name(service_name),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Registry(consul),
	)

	// 初始化服务
	service.Init()

	// 初始化db客户端
	dbproxy.Init(service)

	// 注册服务到注册中心
	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}

}

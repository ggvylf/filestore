package main

import (
	"log"
	"time"

	"github.com/ggvylf/filestore/service/account/handler"
	proto "github.com/ggvylf/filestore/service/account/proto"
	"go-micro.dev/v4"
)

var (
	service_name = "go.micro.service.user"
)

func main() {

	// 创建一个微服务
	service := micro.NewService(
		micro.Name(service_name),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
	)

	// 初始化
	service.Init()

	// 注册服务到注册中心
	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}

}

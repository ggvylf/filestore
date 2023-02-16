package main

import (
	"log"
	"time"

	"github.com/ggvylf/filestore/service/apigw/route"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

var (
	service_name = "go.micro.service.apigw"
)

func startApiService() {
	r := route.Router()
	r.Run(":8888")
}

func startRpcService() {
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

	// 注册服务到注册中心
	if err := service.Run(); err != nil {
		log.Println(err)
	}

}

func main() {
	go startApiService()
	startRpcService()
}

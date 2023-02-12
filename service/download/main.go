package main

import (
	"fmt"
	"time"

	cfg "github.com/ggvylf/filestore/service/download/config"
	dlHandler "github.com/ggvylf/filestore/service/download/handler"
	"github.com/ggvylf/filestore/service/download/route"

	dlProto "github.com/ggvylf/filestore/service/download/proto"

	"go-micro.dev/v4"
)

func startRpcService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
	)
	service.Init()

	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlHandler.Download))
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

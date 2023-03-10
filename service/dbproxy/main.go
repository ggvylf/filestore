package main

import (
	"log"
	"time"

	handler "github.com/ggvylf/filestore/service/dbproxy/handler"

	dbProxy "github.com/ggvylf/filestore/service/dbproxy/proto"

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
		micro.Name("go.micro.service.dbproxy"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),      // 声明超时时间, 避免consul不主动删掉已失去心跳的服务节点
		micro.RegisterInterval(time.Second*5),
		micro.Registry(consul),
	)

	service.Init()

	dbProxy.RegisterDBProxyServiceHandler(service.Server(), new(handler.DBProxy))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

func main() {
	startRpcService()
	// res, err := mapper.FuncCall("/user/UserExist", []interface{}{"haha"}...)
	// log.Printf("error: %+v\n", err)
	// log.Printf("result: %+v\n", res[0].Interface())

	// res, err = mapper.FuncCall("/user/UserExist", []interface{}{"admin"}...)
	// log.Printf("error: %+v\n", err)
	// log.Printf("result: %+v\n", res[0].Interface())
}

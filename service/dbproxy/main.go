package main

import (
	"log"
	"time"

	handler "github.com/ggvylf/filestore/service/dbproxy/handler"

	dbProxy "github.com/ggvylf/filestore/service/dbproxy/proto"

	dbConn "github.com/ggvylf/filestore/service/dbproxy/conn"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

func startRpcService() {
	consul := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Name("go.micro.service.dbproxy"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),      // 声明超时时间, 避免consul不主动删掉已失去心跳的服务节点
		micro.RegisterInterval(time.Second*5),
		micro.Registry(consul),
	)

	service.Init()

	// 初始化db连接
	dbConn.InitDBConn()

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

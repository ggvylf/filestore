package main

import (
	"log"
	"github.com/ggvylf/filestore/service/account/handler"
	proto "github.com/ggvylf/filestore/service/account/handler/proto"
	micro "go-micro.dev/v4"
)

var (
	service_name = "go.micro.service.user"
)

func main() {
	service := micro.NewService(
		micro.Name(service_name)
	)
	service.Init()

	proto.RegisterUserServiceHandler(service.Server(), new(hanlder.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

package main

import (
	"log"

	"github.com/ggvylf/filestore/service/account/handler"
	proto "github.com/ggvylf/filestore/service/account/proto"
	"go-micro.dev/v4"
)

var (
	service_name = "go.micro.service.user"
)

func main() {
	service := micro.NewService(
		micro.Name(service_name),
	)
	service.Init()

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

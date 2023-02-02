package main

import (
	"github.com/fvbock/endless"
	"github.com/ggvylf/filestore/config"
	"github.com/ggvylf/filestore/route"
)

func main() {
	// 加载路由表
	router := route.Route()
	// 优雅退出
	endless.ListenAndServe(config.UploadServiceHost, router)

}

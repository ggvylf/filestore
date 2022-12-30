package main

import (
	"fmt"
	"net/http"

	"github.com/ggvylf/filestore/handler"
)

func main() {
	//上传文件
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)

	//启动web服务
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Printf("failed to listen port,err=%v\n", err.Error())
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/ggvylf/filestore/handler"
)

func main() {
	// 上传文件
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)

	// 下载文件
	http.HandleFunc("/file/download", handler.DownFileHandler)

	// 更新文件
	http.HandleFunc("/file/update", handler.FmUpdateHandler)

	// 删除文件
	http.HandleFunc("/file/delete", handler.FmDeleteHander)

	// 查看指定文件sha1对应的元信息
	http.HandleFunc("/file/meta", handler.GetFileMetaHander)
	http.HandleFunc("/file/meta/all", handler.GetFmListHandler)

	// 启动web服务
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Printf("failed to listen port,err=%v\n", err.Error())
	}
}

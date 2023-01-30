package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ggvylf/filestore/handler"
)

func main() {

	//处理静态资源映射
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

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

	// 秒传
	http.HandleFunc("/file/fastupload", handler.AuthInterceptor(handler.TryFastUploadHandler))

	// 分块上传
	http.HandleFunc("/file/mpupload/init", handler.AuthInterceptor(handler.InitMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", handler.AuthInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", handler.AuthInterceptor(handler.CompleteUploadHandler))

	//获取文件下载的url
	http.HandleFunc("/file/downloadurl", handler.AuthInterceptor(handler.DownloadUrlHandler))

	// 用户注册
	http.HandleFunc("/user/signup", handler.UserSignUpHandler)
	//用户登录
	http.HandleFunc("/user/signin", handler.UserSignInHandler)

	// 用户信息
	// 这里使用了拦截器
	http.HandleFunc("/user/info", handler.AuthInterceptor(handler.UserInfoHandler))

	// 启动web服务
	var srv http.Server
	srv.Addr = ":8888"

	// 优雅退出
	idleConnsClosed := make(chan struct{})
	go func() {

		// 捕获系统信号
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			fmt.Printf("HTTP server Shutdown: %v", err)
		}
		fmt.Println("Closed")
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		fmt.Printf("Failed to listen port,err=%v\n", err.Error())
	}

	// 等待后台goroutine运行结束
	<-idleConnsClosed

}

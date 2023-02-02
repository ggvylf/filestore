package route

import (
	"github.com/ggvylf/filestore/handler"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	// 默认路由
	router := gin.Default()

	// 静态资源
	router.Static("/static", "../../static")

	// 不需要中间件的接口
	// 用户注册
	router.GET("/user/signup", handler.UserSignUpGet)
	router.POST("/user/signup", handler.UserSignUpPost)
	// 用户登录
	router.GET("/user/signin", handler.UserSigninGet)
	router.POST("/user/signin", handler.UserSigninPost)

	// 文件上传成功页面
	router.GET("/file/upload/suc", handler.UploadSucHandler)

	// 加载中间件
	router.Use(handler.AuthInterceptor())

	// 需要中间件的接口

	// 用户信息
	router.POST("/user/info", handler.UserInfoHandler)

	// 上传文件
	router.GET("/file/upload", handler.UploadHandlerGet)
	router.POST("/file/upload", handler.UploadHandlerPost)

	// // 下载文件
	// http.HandleFunc("/file/download", handler.DownFileHandler)

	// // 更新文件
	// http.HandleFunc("/file/update", handler.FmUpdateHandler)

	// // 删除文件
	// http.HandleFunc("/file/delete", handler.FmDeleteHander)

	// 查看指定文件sha1对应的元信息
	router.POST("/file/meta", handler.GetFileMetaHander)
	router.POST("/file/meta/all", handler.GetFmListHandler)

	// // 秒传
	// http.HandleFunc("/file/fastupload", handler.AuthInterceptor(handler.TryFastUploadHandler))

	// // 分块上传
	// http.HandleFunc("/file/mpupload/init", handler.AuthInterceptor(handler.InitMultipartUploadHandler))
	// http.HandleFunc("/file/mpupload/uppart", handler.AuthInterceptor(handler.UploadPartHandler))
	// http.HandleFunc("/file/mpupload/complete", handler.AuthInterceptor(handler.CompleteUploadHandler))

	// //获取文件下载的url
	// http.HandleFunc("/file/downloadurl", handler.AuthInterceptor(handler.DownloadUrlHandler))

	return router
}

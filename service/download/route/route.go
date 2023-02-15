package route

import (
	"github.com/ggvylf/filestore/service/download/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router : 路由表配置
func Router() *gin.Engine {
	// gin framework, 包括Logger, Recovery
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/", "../../../static")

	// // 加入中间件，用于校验token的拦截器(将会从account微服务中验证)
	// router.Use(handler.HTTPInterceptor())

	// 使用gin插件支持跨域请求
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"}, // []string{"http://127.0.0.1:8080"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition", "Authorization"},
		// AllowCredentials: true,
	}))

	// Use之后的所有handler都会经过拦截器进行token校验

	// 文件下载相关接口
	router.GET("/file/download", handler.DownloadHandler)
	router.POST("/file/downloadurl", handler.DownloadURLHandler)

	return router
}

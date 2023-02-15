package route

import (
	"github.com/ggvylf/filestore/service/apigw/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 把静态文件打包到二进制文件中
// type binaryFileSystem struct {
// 	fs http.FileSystem
// }

// func (b *binaryFileSystem) Open(name string) (http.File, error) {
// 	return b.fs.Open(name)
// }

// func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

// 	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
// 		if _, err := b.fs.Open(p); err != nil {
// 			return false
// 		}
// 		return true
// 	}
// 	return false
// }

// func BinaryFileSystem(root string) *binaryFileSystem {
// 	fs := &assetfs.AssetFS{
// 		Asset:     assets.Asset,
// 		AssetDir:  assets.AssetDir,
// 		AssetInfo: assets.AssetInfo,
// 		Prefix:    root,
// 	}
// 	return &binaryFileSystem{
// 		fs,
// 	}
// }

// Router : 网关api路由
func Router() *gin.Engine {
	router := gin.Default()

	// 将静态文件打包到bin文件
	// https://github.com/gin-gonic/examples/blob/master/assets-in-binary/example01/README.md
	// https://github.com/gin-gonic/examples/blob/master/assets-in-binary/example02/README.md
	// router.Use(static.Serve("/static/", BinaryFileSystem("static")))
	router.Static("/static", "../../static")

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"}, // []string{"http://127.0.0.1:8080"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition", "Authorization"},
		// AllowCredentials: true,
	}))

	// 用户注册
	router.GET("/user/signup", handler.UserSignUpGet)
	router.POST("/user/signup", handler.UserSignUpPost)

	// 用户登录
	router.GET("/user/signin", handler.UserSigninGet)
	router.POST("/user/signin", handler.UserSigninPost)

	// 中间件 验证token
	router.Use(handler.AuthInterceptor())

	// 用户信息查询
	router.POST("/user/info", handler.UserInfoHandler)

	// 用户文件列表查询
	router.POST("/file/query", handler.GetFmListHandler)

	// 用户文件修改(重命名)
	router.POST("/file/update", handler.FmUpdateHandler)

	return router
}

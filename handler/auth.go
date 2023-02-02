package handler

import (
	"net/http"

	"github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/util"
	"github.com/gin-gonic/gin"
)

// 拦截器模式
func AuthInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		if len(username) < 3 || !IsTokenValid(username, token) {

			// 验证失败后续不再执行
			c.Abort()
			resp := util.NewRespMsg(
				int(common.StatusTokenInvalid),
				"token无效",
				nil,
			)

			c.JSON(http.StatusOK, resp)

			return

		}

		c.Next()
	}
}

package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ggvylf/filestore/config"
	dblayer "github.com/ggvylf/filestore/db"
	"github.com/ggvylf/filestore/util"
	"github.com/gin-gonic/gin"
)

// 用户注册
func UserSignUpGet(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")

}

func UserSignUpPost(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	// 对用户名和密码做简单的校验
	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "invalid parameter",
			"code": -1,
		})
		return
	}

	// 加密密码
	encpwd := util.Sha1([]byte(passwd + config.PasswordSalt))

	// 用户名 密码写入数据库
	suc := dblayer.UserSignup(username, encpwd)
	if suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "user signup ok!",
			"code": 0,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "user signup failed",
			"code": -2,
		})
	}

}

// 用户登录
func UserSigninGet(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func UserSigninPost(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")
	encpwd := util.Sha1([]byte(passwd + config.PasswordSalt))

	// 从db校验用户名密码
	pwdChecked := dblayer.UserSignin(username, encpwd)
	if !pwdChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "username or password check error!",
			"code": -2,
		})
		return
	}

	// 生成token
	token := GenToken(username)

	// 更新token库
	suc := dblayer.UpdateToken(username, token)
	if !suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "token update failed",
			"code": -2,
		})
		return
	}

	// 登录成功后跳转到主页
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: data{
			Location:      "http://" + c.Request.Host + "/static/view/home.html",
			Username:      username,
			Token:         token,
			DownloadEntry: config.DownloadLBHost,
			UploadEntry:   config.UploadLBHost,
		},
	}

	// c.Data(http.StatusOK, "application/json", resp.JSONBytes())
	c.JSON(http.StatusOK, resp)

}

// 生成token
func GenToken(username string) string {
	// token=md5(usernaem+timestamp+token_salt)+timestamp[:8]
	// len(token)=40
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token_prefix := util.MD5([]byte(username + ts + config.PasswordSalt))
	return token_prefix + ts[:8]
}

type data struct {
	Location      string
	Username      string
	Token         string
	DownloadEntry string
	UploadEntry   string
}

// 用户信息
func UserInfoHandler(c *gin.Context) {

	username := c.Request.FormValue("username")

	// 查询db
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusForbidden, "")
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "ok",
		Data: user,
	}

	c.JSON(http.StatusOK, resp)
}

// token 验证
func IsTokenValid(username, token string) bool {
	// 判断token长度是否是40
	if len(token) < 40 {
		return false
	}

	// 判断token是否过期
	// 判断token是否在db中
	// 对比token

	return true
}

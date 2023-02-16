package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	cmn "github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/config"

	"github.com/ggvylf/filestore/util"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/wrapper/breaker/hystrix"
	ratelimit "github.com/go-micro/plugins/v4/wrapper/ratelimiter/ratelimit"
	ratelimit2 "github.com/juju/ratelimit"
	"go-micro.dev/v4"

	userProto "github.com/ggvylf/filestore/service/account/proto"
	dlProto "github.com/ggvylf/filestore/service/download/proto"
	upProto "github.com/ggvylf/filestore/service/upload/proto"
)

var (
	userCli userProto.UserService
	upCli   upProto.UploadService
	dlCli   dlProto.DownloadService
)

func init() {
	// 配置请求容量及qps
	// 总量1000 qps 100
	bRate := ratelimit2.NewBucketWithRate(100, 1000)

	service := micro.NewService(
		micro.Flags(cmn.CustomFlags...),

		//加入限流功能, false为不等待(超限即返回请求失败)
		micro.WrapClient(ratelimit.NewClientWrapper(bRate, false)),

		// 加入熔断功能, 处理rpc调用失败的情况(cirucuit breaker)
		micro.WrapClient(hystrix.NewClientWrapper()),
	)
	// 初始化， 解析命令行参数等
	service.Init()

	// 创建rpc客户端
	cli := service.Client()
	// tracer, err := tracing.Init("apigw service", "<jaeger-agent-host>")
	// if err != nil {
	// 	log.Println(err.Error())
	// } else {
	// 	cli = client.NewClient(
	// 		client.Wrap(mopentracing.NewClientWrapper(tracer)),
	// 	)
	// }

	// 初始化一个account服务的客户端
	userCli = userProto.NewUserService("go.micro.service.user", cli)
	// 初始化一个upload服务的客户端
	upCli = upProto.NewUploadService("go.micro.service.upload", cli)
	// 初始化一个download服务的客户端
	dlCli = dlProto.NewDownloadService("go.micro.service.download", cli)
}

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

	// 调用account的rpc服务
	resp, err := userCli.Signup(context.TODO(), &userProto.ReqSignup{
		Username: username,
		Password: encpwd,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
	})

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
	resp, err := userCli.Signin(context.TODO(), &userProto.ReqSignin{
		Username: username,
		Password: encpwd,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if resp.Code != cmn.StatusOK {
		c.JSON(200, gin.H{
			"msg":  "登录失败",
			"code": resp.Code,
		})
		return
	}

	// // 动态获取上传入口地址
	// upEntryResp, err := upCli.UploadEntry(context.TODO(), &upProto.ReqEntry{})
	// if err != nil {
	// 	log.Println(err.Error())
	// } else if upEntryResp.Code != cmn.StatusOK {
	// 	log.Println(upEntryResp.Message)
	// }

	// // 动态获取下载入口地址
	// dlEntryResp, err := dlCli.DownloadEntry(context.TODO(), &dlProto.ReqEntry{})
	// if err != nil {
	// 	log.Println(err.Error())
	// } else if dlEntryResp.Code != cmn.StatusOK {
	// 	log.Println(dlEntryResp.Message)
	// }

	// 登录成功后跳转到主页
	cliResp := util.RespMsg{
		Code: int(http.StatusOK),
		Msg:  "登录成功",
		Data: data{
			Location:      "http://" + c.Request.Host + "/static/view/home.html",
			Username:      username,
			Token:         resp.Token,
			DownloadEntry: config.DownloadLBHost,
			UploadEntry:   config.UploadLBHost,
		},
	}

	// c.Data(http.StatusOK, "application/json", resp.JSONBytes())
	c.JSON(http.StatusOK, cliResp)

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

	rpcResp, err := userCli.UserInfo(context.TODO(), &userProto.ReqUserInfo{
		Username: username,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	cliResp := util.RespMsg{
		Code: 0,
		Msg:  "ok",
		Data: gin.H{
			"Username":   username,
			"SignupAt":   rpcResp.SignupAt,
			"LastActive": rpcResp.LastActiveAt,
		},
	}
	c.JSON(http.StatusOK, cliResp)
}

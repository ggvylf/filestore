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
	"github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"

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
	// bRate := ratelimit2.NewBucketWithRate(100, 1000)

	consul := consul.NewRegistry(
		registry.Addrs("127.0.0.1:8500"),
	)

	service := micro.NewService(
		micro.Client(grpc.NewClient()),
		micro.Flags(cmn.CustomFlags...),
		micro.Registry(consul),

	// 	//加入限流功能, false为不等待(超限即返回请求失败)
	// 	micro.WrapClient(ratelimit.NewClientWrapper(bRate, false)),

	// 	// 加入熔断功能, 处理rpc调用失败的情况(cirucuit breaker)
	// 	micro.WrapClient(hystrix.NewClientWrapper()),
	)
	// 初始化， 解析命令行参数等
	service.Init()

	// 创建rpc客户端
	// cli := service.Client()
	// tracer, err := tracing.Init("apigw service", "<jaeger-agent-host>")
	// if err != nil {
	// 	log.Println(err.Error())
	// } else {
	// 	cli = client.NewClient(
	// 		client.Wrap(mopentracing.NewClientWrapper(tracer)),
	// 	)
	// }

	// 初始化一个account服务的客户端
	userCli = userProto.NewUserService("go.micro.service.user", service.Client())
	// 初始化一个upload服务的客户端
	upCli = upProto.NewUploadService("go.micro.service.upload", service.Client())
	// 初始化一个download服务的客户端
	dlCli = dlProto.NewDownloadService("go.micro.service.download", service.Client())
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

	// 调用account的rpc服务
	resp, err := userCli.Signup(context.TODO(), &userProto.ReqSignup{
		Username: username,
		Password: passwd,
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

	// 对用户名和密码做简单的校验
	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "invalid parameter",
			"code": cmn.StatusParamInvalid,
		})
		return
	}

	// rpc调用登录服务
	resp, err := userCli.Signin(context.TODO(), &userProto.ReqSignin{
		Username: username,
		Password: passwd,
	})

	if err != nil || resp.Code != cmn.StatusOK {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg":  resp.Message,
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

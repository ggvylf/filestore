package handler

import (
	"context"
	"log"

	"github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/config"

	dbCli "github.com/ggvylf/filestore/service/dbproxy/client"

	proto "github.com/ggvylf/filestore/service/account/proto"
	"github.com/ggvylf/filestore/util"
)

type User struct{}

// 用户注册
func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, resp *proto.RespSignup) error {

	username := req.Username

	passwd := req.Password

	// 对用户名和密码做简单的校验
	if len(username) < 3 || len(passwd) < 5 {
		resp.Code = common.StatusParamInvalid
		resp.Message = "用户参数非法"
		return nil
	}

	// 加密密码
	encpwd := util.Sha1([]byte(passwd + config.PasswordSalt))

	// 用户名 密码写入数据库
	res, err := dbCli.UserSignup(username, encpwd)
	if err != nil || !res.Suc {
		log.Println(err.Error())
		log.Println(res.Msg)
		resp.Code = common.StatusRegisterFailed
		resp.Message = "用户注册失败"

	}
	resp.Code = common.StatusOK
	resp.Message = "用户注册成功"
	return nil

}

func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, resp *proto.RespSignin) error {

	username := req.Username
	passwd := req.Password
	encpwd := util.Sha1([]byte(passwd + config.PasswordSalt))

	// 从db校验用户名密码
	dbResp, err := dbCli.UserSignin(username, encpwd)
	if err != nil || !dbResp.Suc {
		log.Println(err.Error())
		log.Println(dbResp.Msg)
		resp.Code = common.StatusLoginFailed
		resp.Message = "校验用户名密码失败"
		return nil
	}

	// 生成token
	token := util.GenToken(username)

	// 更新token库
	res, err := dbCli.UpdateToken(username, token)
	if err != nil || !res.Suc {
		log.Println(err.Error())
		log.Println(dbResp.Msg)
		resp.Code = common.StatusServerError
		resp.Message = "更新token失败"
		return nil
	}

	resp.Code = common.StatusOK
	resp.Token = token

	return nil
}

// 从tbl_user表中查询用户信息
func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, resp *proto.RespUserInfo) error {

	// 查询db
	res, err := dbCli.GetUserInfo(req.Username)
	if err != nil || !res.Suc {
		log.Println(err.Error())
		resp.Code = common.StatusUserNotExists
		resp.Message = "用户不存在"
		return nil
	}

	// 序列化用户信息
	user := dbCli.ToTableUser(res.Data)

	resp.Code = common.StatusOK
	resp.Username = user.Username
	resp.SignupAt = user.SignupAt
	resp.LastActiveAt = user.LastActiveAt
	resp.Status = int64(user.Status)
	// TODO: 需增加接口支持完善用户信息(email/phone等)
	resp.Email = user.Email
	resp.Phone = user.Phone

	return nil
}

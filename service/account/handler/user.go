package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/config"
	dbcli "github.com/ggvylf/filestore/service/dbproxy"

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
		resp.Message = "invalid parameter"
		return nil
	}

	// 加密密码
	encpwd := util.Sha1([]byte(passwd + config.PasswordSalt))

	// 用户名 密码写入数据库
	suc := dbcli.UserSignup(username, encpwd)
	if suc {
		resp.Code = common.StatusOK
		resp.Message = "user signup suc"

	} else {
		resp.Code = common.StatusRegisterFailed
		resp.Message = "user signup failed"
	}

	return nil

}

func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, resp *proto.RespSignin) error {

	username := req.Username
	passwd := req.Password
	encpwd := util.Sha1([]byte(passwd + config.PasswordSalt))

	// 从db校验用户名密码
	pwdChecked := dbcli.UserSignin(username, encpwd)
	if !pwdChecked {
		resp.Code = common.StatusLoginFailed
		return nil
	}

	// 生成token
	token := GenToken(username)

	// 更新token库
	suc := dblayer.UpdateToken(username, token)
	if !suc {

		resp.Code = common.StatusServerError
		return nil
	}

	resp.Code = common.StatusOK
	resp.Token = token
	return nil
}

// 生成token
func GenToken(username string) string {
	// token=md5(usernaem+timestamp+token_salt)+timestamp[:8]
	// len(token)=40
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token_prefix := util.MD5([]byte(username + ts + config.PasswordSalt))
	return token_prefix + ts[:8]
}

// 从tbl_user表中查询用户信息
func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, resp *proto.RespUserInfo) error {

	// 查询db
	dbresp, err := dbcli.GetUserInfo(req.Username)
	if err != nil {
		resp.Code = common.StatusUserNotExists
		resp.Message = "用户不存在"
		return nil
	}

	//
	user := dbcli.ToTableUser(dbresp.Data)

	resp.Code = common.StatusOK
	resp.Username = user.Username
	resp.SignupAt = user.SignupAt
	resp.LastActiveAt = user.LastActiveAt
	resp.Status = int32(user.Status)
	// TODO: 需增加接口支持完善用户信息(email/phone等)
	resp.Email = user.Email
	resp.Phone = user.Phone

	return nil
}

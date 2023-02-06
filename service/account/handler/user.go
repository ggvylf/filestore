package handler

import (
	"context"

	"github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/config"
	dblayer "github.com/ggvylf/filestore/db"

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
	suc := dblayer.UserSignup(username, encpwd)
	if suc {
		resp.Code = common.StatusOK
		resp.Message = "user signup suc"

	} else {
		resp.Code = common.StatusRegisterFailed
		resp.Message = "user signup failed"
	}

	return nil

}

package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	dblayer "github.com/ggvylf/filestore/db"
	"github.com/ggvylf/filestore/util"
)

var (
	pwd_salt   = "mysalt"
	token_salt = "mytoken"
)

// 用户注册
func UserSignUpHandler(w http.ResponseWriter, r *http.Request) {
	//Get 方式返回注册页面
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signup.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}

	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	// 对用户名和密码做简单的校验
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}

	// 加密密码
	encpwd := util.Sha1([]byte(passwd + pwd_salt))

	// 用户名 密码写入数据库
	suc := dblayer.UserSignup(username, encpwd)
	if suc {
		w.Write([]byte("user signup success"))
	} else {
		w.Write([]byte("user signup failed"))
	}

}

// 用户登录
func UserSignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	encpwd := util.Sha1([]byte(passwd + pwd_salt))

	// 从db校验用户名密码
	pwdChecked := dblayer.UserSignin(username, encpwd)
	if !pwdChecked {
		w.Write([]byte("user checked failed"))
		return
	}

	// 生成token
	token := GenToken(username)
	// 更新token库
	suc := dblayer.UpdateToken(username, token)
	if !suc {
		w.Write([]byte("User login failed"))
		return
	}

	// 登录成功后跳转到主页
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))

}

// 生成token
func GenToken(username string) string {
	// token=md5(usernaem+timestamp+token_salt)+timestamp[:8]
	// len(token)=40
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token_prefix := util.MD5([]byte(username + ts + token_salt))
	return token_prefix + ts[:8]
}

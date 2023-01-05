package db

import (
	"fmt"

	mydb "github.com/ggvylf/filestore/db/mysql"
)

// 用户注册
func UserSignup(username, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user (user_name,user_pwd) values(?,?)")
	if err != nil {
		fmt.Println("Failed to insert new user,err=" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)

	if err != nil {
		fmt.Println("Failed to insert new user,err=" + err.Error())
		return false
	}

	// 如果影响行数大于0 表示插入成功
	rf, err := ret.RowsAffected()
	if err == nil && rf > 0 {
		return true
	}
	return false

}

// 用户登录
func UserSignin(username, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println("User not exist,err=" + err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("User not exist,err=" + err.Error())
		return false
	} else if rows == nil {
		fmt.Println("User not exist,err=" + err.Error())
		return false
	}

	// 对返回数据写入一个map
	pRows := mydb.ParseRows(rows)

	// 取map的第一个值中的user_pwd字段
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// 更新用户token
func UpdateToken(username, token string) bool {
	stmt, err := mydb.DBConn().Prepare("insert into tbl_user_token(user_name,user_token) values(?,?)")
	if err != nil {
		fmt.Println("Update User Token Failed,err=" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, token)
	if err != nil {
		fmt.Println("Update User Token Failed,err=" + err.Error())
		return false
	}

	// 如果影响行数大于0 表示插入成功
	rf, err := ret.RowsAffected()
	if err == nil && rf > 0 {
		return true
	}
	return false
}

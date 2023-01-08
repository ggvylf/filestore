package db

import (
	"fmt"
	"time"

	mydb "github.com/ggvylf/filestore/db/mysql"
)

// 跟tbl_user_file表中字段保持一致
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    string
	UploadAt    string
	LastUpdated string
}

func UpdateUserFile(usernaem, filehash, filename, filesize string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user_file (user_name,file_sha1,file_name,file_size,upload_at) values(?,?,?,?,?)")
	if err != nil {
		fmt.Println("failed to conn db,err" + err.Error())
		return false
	}
	defer stmt.Close()

	// ret, err := stmt.Exec(usernaem, filehash, filename, filesize, time.Now())
	_, err = stmt.Exec(usernaem, filehash, filename, filesize, time.Now())

	if err != nil {
		fmt.Println("update tbl_user_file failed,err" + err.Error())
		return false
	}

	// rf, err := ret.RowsAffected()
	// if rf <= 0 {
	// 	fmt.Println("update tbl_user_file failed,err" + err.Error())
	// 	return false
	// }

	return true
}

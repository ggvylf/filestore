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

// 更新tbl_user_file
func UpdateUserFile(username, filehash, filename string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user_file (user_name,file_sha1,file_name,file_size,upload_at) values(?,?,?,?,?)")
	if err != nil {
		fmt.Println("failed to conn db,err" + err.Error())
		return false
	}
	defer stmt.Close()

	// ret, err := stmt.Exec(usernaem, filehash, filename, filesize, time.Now())
	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())

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

// 批量从tbl_user_file中获取原信息
func GetUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_name,file_size,upload_at,last_update  from tbl_user_file where user_name=? limit ?")
	if err != nil {
		fmt.Println("failed to conn db,err" + err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Println("select tbl_user_file failed,err" + err.Error())
		return nil, err
	}

	var UserFiles []UserFile
	for rows.Next() {
		uf := UserFile{}
		err = rows.Scan(&uf.FileHash, &uf.FileName, &uf.FileSize, &uf.UploadAt, &uf.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		UserFiles = append(UserFiles, uf)

	}

	return UserFiles, nil

}

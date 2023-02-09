package orm

import (
	"fmt"
	"log"
	"time"

	mydb "github.com/ggvylf/filestore/service/dbproxy/conn"
)

// 从tbl_user_file中获取fmlist
func QueryUserFileMetas(username string, limit int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_name,file_size,upload_at,last_update  from tbl_user_file where user_name=? limit ?")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	var UserFiles []TableUserFile
	for rows.Next() {
		uf := TableUserFile{}
		err = rows.Scan(&uf.FileHash, &uf.FileName, &uf.FileSize, &uf.UploadAt, &uf.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		UserFiles = append(UserFiles, uf)

	}

	res.Suc = true
	res.Data = UserFiles
	return

}

// 从tbl_user_file中获取fm
func QueryUserFileMeta(username, filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_name,file_size,upload_at,last_update  from tbl_user_file where user_name=? and file_sha1=?")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	var UserFiles []TableUserFile
	for rows.Next() {
		uf := TableUserFile{}
		err = rows.Scan(&uf.FileHash, &uf.FileName, &uf.FileSize, &uf.UploadAt, &uf.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		UserFiles = append(UserFiles, uf)

	}

	res.Suc = true
	res.Data = UserFiles
	return

}

// 重命名文件，修改tbl_user_file
func RenameFileName(username, filehash, newfilename string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare("update tbl_user_file set file_name=? where file_sha1=? and user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(newfilename, filehash, username)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	res.Suc = true
	return

}

// 标记删除文件 tbl_user_file
func DeleteUserFile(username, filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_user_file set status=2 where user_name=? and file_sha1=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// 把fm插入tbl_user_file
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_at`) values (?,?,?,?,?)")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

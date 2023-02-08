package orm

import (
	"fmt"
	"log"

	mydb "github.com/ggvylf/filestore/service/dbproxy/conn"
)

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

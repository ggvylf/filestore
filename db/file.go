package db

import (
	"database/sql"
	"fmt"

	mydb "github.com/ggvylf/filestore/db/mysql"
)

// 插入文件信息到tbl_file
func InsertFmDb(filehash, filename, fileaddr string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")

	if err != nil {
		fmt.Println("Faled to prepare statement,err=", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// 根据影响行数判断insert是否成功
	rf, err := ret.RowsAffected()

	if err == nil {
		if rf <= 0 {
			// 当前数据已存在
			fmt.Printf("File with hash=%s has been existed", filehash)

		}

		return true
	}

	return false

}

// 定义一个结构体 用来存放从db中取出的字段
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// 从tbl_file获取fm的元数据信息
func GetFmDb(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1,file_name,file_addr,file_size from tbl_file where file_sha1 = ? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := TableFile{}

	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileAddr, &tfile.FileSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &tfile, nil
}

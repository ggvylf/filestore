package meta

import (
	"os"

	dblayer "github.com/ggvylf/filestore/db"
)

// 文件的元信息
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UpoadAt  string
}

// 文件metatdata列表 key是FileSha1
var fmList map[string]FileMeta

func init() {
	// 初始化列表
	fmList = make(map[string]FileMeta)
}

// 获取fmlist
func GetFmList() map[string]FileMeta {
	return fmList
}

// 更新fileMetas列表
func UploadFmList(fm FileMeta) {
	fmList[fm.FileSha1] = fm
}

// 获取fileMetas中的FileMeta对象
func GetFm(filehash string) FileMeta {
	return fmList[filehash]
}

// 更新fm元数据到db
func UpdateFmDb(fm FileMeta) bool {
	return dblayer.InsertFmDb(fm.FileSha1, fm.FileName, fm.Location, fm.FileSize)
}

// 从tbl_file中获取fm元信息
func GetFmDb(filehash string) (*FileMeta, error) {
	tfile, err := dblayer.GetFmDb(filehash)
	if err != nil || tfile == nil {
		return nil, err
	}
	fm := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return &fm, nil
}

// 从fmList中删除fm
func DeleteFm(filehash string) {
	delete(fmList, filehash)
}

// 从磁盘删除文件
func DeleteFile(location string) bool {
	err := os.Remove(location)
	if err != nil {
		return false
	}
	return true
}

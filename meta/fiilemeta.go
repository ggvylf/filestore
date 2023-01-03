package meta

import "os"

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
func UploadFmList(meta FileMeta) {
	fmList[meta.FileSha1] = meta
}

// 获取fileMetas中的FileMeta对象
func GetFm(sha1 string) FileMeta {
	return fmList[sha1]
}

// 从fmList中删除fm
func DeleteFm(sha1 string) {
	delete(fmList, sha1)
}

// 从磁盘删除文件
func DeleteFile(location string) bool {
	err := os.Remove(location)
	if err != nil {
		return false
	}
	return true
}

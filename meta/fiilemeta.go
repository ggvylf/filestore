package meta

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

// 更新fileMetas列表
func UploadFmList(meta FileMeta) {
	fmList[meta.FileSha1] = meta
}

// 获取fileMetas中的FileMeta对象
func GetFm(sha1 string) FileMeta {

	return fmList[sha1]
}

func GetFmList() map[string]FileMeta {
	return fmList
}

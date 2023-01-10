package handler

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	rPool "github.com/ggvylf/filestore/cache/redis"
	dblayer "github.com/ggvylf/filestore/db"
	"github.com/ggvylf/filestore/util"
	"github.com/gomodule/redigo/redis"
)

// 分块上传信息结构体
type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	// 上传的id 不重复
	UploadID string
	// 分块的大小
	ChunkSize int
	// 分块数量
	ChunkCount int
}

// 初始化分块上传
// 用redis保存相关信息
func InitMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// 获取redis连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 初始化分块信息
	mpInfo := MultipartUploadInfo{
		FileHash:  filehash,
		FileSize:  filesize,
		UploadID:  username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize: 5 * 1024 * 1024, //5M,

		// 分块的个数
		ChunkCount: int(math.Ceil(float64(filesize)) / (5 * 1024 * 1024)),
	}

	// 写入redis

	// 期望的分块数量
	rConn.Do("HSET", "MP_"+mpInfo.UploadID, "chunkcount", mpInfo.ChunkCount)

	rConn.Do("HSET", "MP_"+mpInfo.UploadID, "filehash", mpInfo.FileHash)
	rConn.Do("HSET", "MP_"+mpInfo.UploadID, "filesize", mpInfo.FileSize)

	// 返回响应信息
	w.Write(util.NewRespMsg(0, "ok", mpInfo).JSONBytes())

}

// 分块上传接口
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// 获取redis连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 创建文件句柄
	fpath := "/tmp" + uploadID + "/" + chunkIndex
	os.MkdirAll(fpath, 0644)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "upload part failed", nil).JSONBytes())
		return

	}

	defer fd.Close()

	// 写入文件
	buf := make([]byte, 1024*1024) //1MB

	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 更新redis中的分块记录
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回处理结果
	w.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}

// 分块合并接口
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	r.ParseForm()
	uploadid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")
	chunkindex := r.Form.Get("index")

	// 获取redis连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 通过uploadid查询redis 判断分块是否全部上传完成
	// redis里查不到记录 直接返回
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
		return

	}

	// 验证redis中的分块记录

	// 期望的数量
	total := 0

	// 实际分块的数量
	chunkcount := 0


	//
	chunkindex=

	// data中的格式是k1 v1 k2 v2
	for i := 0; i < len(data); i += 2 {

		k := string(data[i].(byte))
		v := string(data[i+1].(byte))

		// 获取期望数量
		if k == "chunkcount" {
			total, _ = strconv.Atoi(v)

			// 实际分块的数量
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkcount++
		}

	}

	if total != chunkcount {
		w.Write(util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
		return
	}

	// 合并分块

	// 找到要合并的文件
	fpath := "/tmp" + uploadid + "/" + chunkIndex
	
	// 合并操作

	fileaddr

	// 更新tbl_file和tbl_user_file
	fsize, _ := strconv.Atoi(filesize)
	dblayer.InsertFmDb(filehash, filename, fileaddr, int64(fsize))
	dblayer.UpdateUserFile(username, filehash, filename, int64(fsize))

	// 响应处理结果
	w.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}

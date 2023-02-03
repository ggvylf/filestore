package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	rPool "github.com/ggvylf/filestore/cache/redis"
	"github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/config"
	dblayer "github.com/ggvylf/filestore/db"
	"github.com/ggvylf/filestore/mq"
	"github.com/ggvylf/filestore/util"
	"github.com/gin-gonic/gin"
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
func InitMultipartUploadHandler(c *gin.Context) {
	// 解析参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.Data(http.StatusOK, "", util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
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
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// 写入redis

	// 期望的分块数量
	rConn.Do("HSET", "MP_"+mpInfo.UploadID, "chunkcount", mpInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+mpInfo.UploadID, "filehash", mpInfo.FileHash)
	rConn.Do("HSET", "MP_"+mpInfo.UploadID, "filesize", mpInfo.FileSize)

	// example:
	// 127.0.0.1:6379> HGETALL MP_admin17397e0729ee9c44
	// 1) "chunkcount"
	// 2) "28"
	// 3) "filehash"
	// 4) "fe1d6ccb2544698b5c567411306e659de0fe922d"
	// 5) "filesize"
	// 6) "148883574"

	// 返回响应信息
	c.Data(http.StatusOK, "", util.NewRespMsg(0, "ok", mpInfo).JSONBytes())

}

// 分块上传接口
func UploadPartHandler(c *gin.Context) {
	// 解析参数
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 获取redis连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 创建文件句柄
	fpath := "/tmp" + "/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.Data(http.StatusOK, "", util.NewRespMsg(-1, "create part failed", nil).JSONBytes())
		return

	}

	defer fd.Close()

	// 写入文件
	buf := make([]byte, 1024*1024) //1MB

	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 更新redis中的分块记录
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回处理结果
	c.Data(http.StatusOK, "", util.NewRespMsg(0, "ok", nil).JSONBytes())
}

// 分块合并接口
func CompleteUploadHandler(c *gin.Context) {
	// 解析参数

	uploadid := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize := c.Request.FormValue("filesize")
	filename := c.Request.FormValue("filename")

	// 获取redis连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 通过uploadid查询redis 判断分块是否全部上传完成
	// redis里查不到记录 直接返回
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadid))
	if err != nil {
		c.Data(http.StatusOK, "", util.NewRespMsg(-1, "Complete upload failed", nil).JSONBytes())
		return

	}

	// 验证redis中的分块记录

	// 期望的数量
	total := 0

	// 实际分块的数量
	chunkcount := 0

	// data中的格式是k1 v1 k2 v2
	for i := 0; i < len(data); i += 2 {

		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))

		// 获取期望数量
		if k == "chunkcount" {
			total, _ = strconv.Atoi(v)

			// 实际分块的数量
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkcount++
		}

	}

	if total != chunkcount {
		c.Data(http.StatusOK, "", util.NewRespMsg(-1, "chunkcount check failed", chunkcount).JSONBytes())
		return
	}

	// 合并分块
	fpath := "/tmp" + "/" + uploadid + "/"
	_, fname := path.Split(filename)
	fileaddr := fmt.Sprintf("/tmp/" + fname)

	fd, _ := os.OpenFile(fileaddr, os.O_CREATE|os.O_WRONLY, 0644)
	defer fd.Close()

	files, _ := filepath.Glob(fpath + "*")
	// fmt.Println(files)
	filessorted, err := util.FileSortForStringWithNum(files)
	if err != nil {
		fmt.Println("files sort failed,err=", err)
	}
	// fmt.Println(filessorted)

	for _, f := range filessorted {
		// 排除目标文件在同目录下
		if filepath.Base(f) == fname {
			break
		}
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("read part file err=", err)
		}
		_, err = fd.Write(data)
		if err != nil {
			fmt.Println("write part file err=", err)
		}
	}

	fmt.Println("complete file suc")

	// 更新tbl_file和tbl_user_file
	fsize, _ := strconv.Atoi(filesize)

	name := filepath.Base(filename)
	dblayer.InsertFmDb(filehash, name, fileaddr, int64(fsize))
	dblayer.UpdateUserFile(username, filehash, name, int64(fsize))

	// 拼接mq信息

	bucket := "userfile"
	ossName := "/minio" + "/" + filehash

	msgdata := mq.TransferData{
		FileHash:      filehash,
		CurLocation:   fileaddr,
		DestLocation:  ossName,
		DestStoreType: common.StoreOSS,
		FileSize:      int64(fsize),
		Bucket:        bucket,
	}
	pubData, _ := json.Marshal(msgdata)

	// 推送到mq
	suc := mq.Publush(config.TransExchangeName, config.TransOSSRoutingKey, pubData)
	if !suc {
		fmt.Println("send to mq failed")
		return
	}

	// 响应处理结果
	c.Data(http.StatusOK, "", util.NewRespMsg(0, "ok", nil).JSONBytes())

}

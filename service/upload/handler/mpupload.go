package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	rPool "github.com/ggvylf/filestore/cache/redis"
	"github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/config"
	"github.com/ggvylf/filestore/mq"
	dbcli "github.com/ggvylf/filestore/service/dbproxy/client"
	"github.com/ggvylf/filestore/util"
	"github.com/gin-gonic/gin"
)

// MultipartUploadInfo : 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// InitialMultipartUploadHandler : 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {
	// 1. 解析用户请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -1,
				"msg":  "params invalid",
			})
		return
	}

	// 2. 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// 4. 将初始化信息写入到redis缓存
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// 5. 将响应初始化数据返回到客户端
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "OK",
			"data": upInfo,
		})
}

// UploadPartHandler : 上传文件分块
func UploadPartHandler(c *gin.Context) {
	// 解析参数
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 获取redis连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 创建文件句柄
	fpath := config.TempPartRootDir + "/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		c.Data(http.StatusOK, "", util.NewRespMsg(-1, "create part dir failed", nil).JSONBytes())
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

// CompleteUploadHandler : 通知上传合并
func CompleteUploadHandler(c *gin.Context) {

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
	fpath := config.TempPartRootDir + "/" + uploadid + "/"
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

	fm := dbcli.FileMeta{
		FileSha1: filehash,
		FileName: name,
		FileSize: int64(fsize),
		Location: fileaddr,
	}

	// 更新tbl_file
	_, ferr := dbcli.OnFileUploadFinished(fm)

	// 更新tbl_user_file
	_, uferr := dbcli.OnUserFileUploadFinished(username, fm)

	if ferr != nil || uferr != nil {
		log.Println(err)
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": -2,
				"msg":  "数据更新失败",
				"data": nil,
			})
		return
	}

	// 拼接mq信息
	bucket := config.OSSBucket
	ossName := config.OSSRootDir + "/" + filehash

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

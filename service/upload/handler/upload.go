package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/config"
	cfg "github.com/ggvylf/filestore/config"
	"github.com/ggvylf/filestore/mq"
	dbcli "github.com/ggvylf/filestore/service/dbproxy/client"
	"github.com/ggvylf/filestore/service/dbproxy/orm"
	store "github.com/ggvylf/filestore/store/minio"
	"github.com/ggvylf/filestore/util"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// DoUploadHandler ： 处理文件上传
func DoUploadHandler(c *gin.Context) {

	// 错误码
	// 0 没错误
	// -1 获取form错误
	// -2 把form内容转换成[]byte
	// -3 打开本地文件描述符错误
	// -4 文件写入本地错误
	// -5 写入oss或者是ceph错误
	// -6 操作db错误

	errCode := 0

	// 函数关闭的时候返回结果
	defer func() {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传成功",
			})
		}
	}()

	//从form表单中获取文件
	file, err := c.FormFile("file")

	if err != nil {
		log.Printf("failed to get data,err=%v\n", err.Error())
		errCode = -1
		return
	}

	// 初始化FileMeta
	fm := dbcli.FileMeta{
		FileName: file.Filename,
		Location: cfg.TempLocalRootDir + "/" + file.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		FileSize: file.Size,
	}

	// //新建一个本地文件的fd
	// newfile, err := os.Create(fm.Location)
	// if err != nil {
	// 	fmt.Printf("Failed to create file,err=%v\n", err.Error())
	// 	return
	// }
	// defer newfile.Close()

	// // 复制文件 同时可以获取文件大小
	// fm.FileSize, err = io.Copy(newfile, file)
	// if err != nil {
	// 	fmt.Printf("Failed to write file,err=%v\n", err.Error())
	// 	return
	// }

	// 复制文件到本地
	err = c.SaveUploadedFile(file, fm.Location)
	if err != nil {
		log.Printf("Failed to write file,err=%v\n", err.Error())
		errCode = -4
		return
	}

	// 重置newfile的偏移量到文件头部
	// newfile.Seek(0, 0)

	// 计算上传文件的sha1
	fm.FileSha1 = util.FileSha1(fm.Location)

	//把文件写入对象存储
	data, _ := os.Open(fm.Location)
	ctx := context.Background()
	mc := store.GetMC()
	bucket := cfg.OSSBucket
	ossName := cfg.OSSRootDir + "/" + fm.FileSha1
	// path := "/userfile" + ossName

	// _, err = mc.PutObject(ctx, bucket, ossName, data, fm.FileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	// if err != nil {
	// 	fmt.Println("upload file to oss failed,err=", err)
	// 	return
	// }
	// fm.Location = path

	// 写入文件
	// 根据配置文件确定写入的类型，以及是否需要用到mq来异步

	// 写入ceph
	if cfg.CurrentStoreType == common.StoreCeph {
		mc.PutObject(ctx, bucket, ossName, data, fm.FileSize, minio.PutObjectOptions{})
		// 写入minio
	} else if cfg.CurrentStoreType == common.StoreOSS {

		// 写入oss是同步还是异步
		if !cfg.AsyncTransferEnable {

			// 同步写入oss
			mc.PutObject(ctx, bucket, ossName, data, fm.FileSize, minio.PutObjectOptions{})

			// 异步写入mq
		} else {

			// 拼接mq信息
			data := mq.TransferData{
				FileHash:      fm.FileSha1,
				CurLocation:   fm.Location,
				DestLocation:  ossName,
				DestStoreType: common.StoreOSS,
				FileSize:      file.Size,
				Bucket:        bucket,
			}
			pubData, _ := json.Marshal(data)

			// 推送到mq
			suc := mq.Publush(config.TransExchangeName, config.TransOSSRoutingKey, pubData)
			if !suc {
				log.Println("send to mq failed")
				errCode = -5
				return
			}

		}

		// append到元信息队列中
		// meta.UploadFmList(fm)

		// 更新元数据到文件表 tbl_file
		_, err = dbcli.OnFileUploadFinished(fm)
		if err != nil {
			log.Println(err.Error())
			errCode = -6
			return
		}

		// 更新信息到用户文件表 tbl_user_file
		// BUG(myself): 前端页面FORM表单里没有username
		username := c.Request.FormValue("username")
		upResp, suc := dbcli.OnUserFileUploadFinished(username, fm)
		if suc != nil || upResp.Suc {
			c.Data(http.StatusOK, "text/plain", []byte("update db failed!"))
			errCode = -6
		}

		errCode = 0

		// 302重定向到上传成功页面
		// c.Redirect(http.StatusFound, "/file/upload/suc")
	}
}

// TryFastUploadHandler : 尝试秒传接口
func TryFastUploadHandler(c *gin.Context) {

	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	// filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	// 2. 从文件表中查询相同hash的文件记录
	fileMetaResp, err := dbcli.GetFileMeta(filehash)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if !fileMetaResp.Suc {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}

	// 4. 上传过则将文件信息写入用户文件表， 返回成功
	fmeta := dbcli.TableFileToFileMeta(fileMetaResp.Data.(orm.TableFile))
	fmeta.FileName = filename
	upRes, err := dbcli.OnUserFileUploadFinished(username, fmeta)
	if err == nil && upRes.Suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
	return
}

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ggvylf/filestore/common"
	"github.com/ggvylf/filestore/config"
	dblayer "github.com/ggvylf/filestore/db"
	"github.com/ggvylf/filestore/meta"
	"github.com/ggvylf/filestore/mq"
	userProto "github.com/ggvylf/filestore/service/account/proto"
	store "github.com/ggvylf/filestore/store/minio"
	"github.com/ggvylf/filestore/util"
	"github.com/gin-gonic/gin"
)

// 上传文件并保存到本地
func UploadHandlerGet(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/upload.html")
}

func UploadHandlerPost(c *gin.Context) {
	//从form表单中获取文件
	file, err := c.FormFile("file")

	if err != nil {
		fmt.Printf("failed to get data,err=%v\n", err.Error())
		return
	}

	// 初始化FileMeta
	fm := meta.FileMeta{
		FileName: file.Filename,
		Location: "/tmp/" + file.Filename,
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
		fmt.Printf("Failed to write file,err=%v\n", err.Error())
		return
	}

	// // 重置newfile的偏移量到文件头部
	// newfile.Seek(0, 0)

	// 计算上传文件的sha1
	fm.FileSha1 = util.FileSha1(fm.Location)

	//把文件写入对象存储
	// data, _ := os.Open(fm.Location)
	// ctx := context.Background()
	// mc := store.GetMC()
	bucket := "userfile"
	ossName := "/minio" + "/" + fm.FileSha1
	// path := "/userfile" + ossName

	// _, err = mc.PutObject(ctx, bucket, ossName, data, fm.FileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	// if err != nil {
	// 	fmt.Println("upload file to oss failed,err=", err)
	// 	return
	// }
	// fm.Location = path

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
		fmt.Println("send to mq failed")
		return
	}

	// append到元信息队列中
	// meta.UploadFmList(fm)

	// 更新元数据到文件表 tbl_file
	_ = meta.UpdateFmDb(fm)

	// 更新信息到用户文件表 tbl_user_file
	// BUG(myself): 前端页面FORM表单里没有username
	username := c.Request.FormValue("username")

	suc = dblayer.UpdateUserFile(username, fm.FileSha1, fm.FileName, fm.FileSize)
	if !suc {
		c.Data(http.StatusOK, "text/plain", []byte("update db failed!"))
		return
	}

	// 302重定向到上传成功页面
	c.Redirect(http.StatusFound, "/file/upload/suc")
}

// 上传文件成功
func UploadSucHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("update ok"))
}

// 返回元信息列表
// 老的方法是从内存中获取
// 新的方法从tbl_user_file中获取
func GetFmListHandler(c *gin.Context) {
	// fmList := meta.GetFmList()

	username := c.Request.FormValue("username")
	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))

	rpcResp, err := userCli.UserFilesList(context.TODO(), &userProto.ReqUserFile{
		Username: username,
		Limit:    int64(limitCnt),
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}

	c.Data(http.StatusOK, "application/json", rpcResp.FileData)

}

// 返回指定sha1的fm对象
func GetFileMetaHander(c *gin.Context) {

	filehash := c.Request.FormValue("filehash")
	// fm := meta.GetFm(filehash)

	// 从db中获取fm
	fm, err := meta.GetFmDb(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
		return
	}
	c.JSON(http.StatusOK, fm)

}

// // 从本地文件下载
// // 已废弃
// func DownFileHandler(c *gin.Context) {

// 	filehash := c.Request.FormValue("filehash")

// 	fm := meta.GetFmDb(filehash)

// 	//小文件可以 大文件性能不行 这两种方法是等价的
// 	//data, err := os.ReadFile(fm.Location)
// 	//data, err := ioutil.ReadFile(fm.FileName)

// 	// 先打开文件句柄再读取
// 	fd, err := os.Open(fm.FileName)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, "")
// 		return
// 	}
// 	defer fd.Close()

// 	data, err := ioutil.ReadAll(fd)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, "")

// 		return
// 	}

// 	c.Writer.Header().Add("Content-Type", "application/octect")
// 	// 避免中文文件名乱码
// 	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(fm.FileName)))

// 	c.Data(http.StatusOK, "", data)

// }

// 更新fm的filename
func FmUpdateHandler(c *gin.Context) {
	fileHash := c.Request.FormValue("filehash")
	opType := c.Request.FormValue("op")
	newFileName := c.Request.FormValue("filename")
	username := c.Request.FormValue("username")

	// 判断optype
	if opType != "0" {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}

	// 重命名文件
	rpcResp, err := userCli.UserFileRename(context.TODO(), &userProto.ReqUserFileRename{
		Username:    username,
		Filehash:    fileHash,
		NewFileName: newFileName,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}

	c.JSON(http.StatusOK, rpcResp.FileData)

}

// // 删除fm和文件
// func FmDeleteHander(c *gin.Context) {
// 	fileHash := c.Request.FormValue("filehash")

// 	fm := meta.GetFmDb(fileHash)

// 	// 这里注意要先删除文件 再删除元信息
// 	ok := meta.DeleteFile(fm.Location)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, []byte("file not exists"))
// 		return
// 	}

// 	meta.DeleteFm(fileHash)
// 	c.JSON(http.StatusOK, []byte("delete ok"))

// }

// 尝试秒传接口
// 秒传
// 1. 判断文件是否有记录在tbl_file中，
// 2. 如果有记录，不用上传，直接更新tbl_user_file信息
// 3. 如果没有记录，走/file/upload接口
func TryFastUploadHandler(c *gin.Context) {

	// 解析参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	// 查询tbl_file中相同filehash
	fm, err := meta.GetFmDb(filehash)

	// 判断文件是否存在
	if fm == nil || err != nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，使用普通上传接口",
		}
		c.JSON(http.StatusOK, resp)

		return
	}

	// 更新tbl_user_file
	suc := dblayer.UpdateUserFile(username, filehash, filename, int64(filesize))
	if !suc {
		resp := util.RespMsg{
			Code: -2,
			Msg:  "秒传失败，请稍后重试",
		}
		c.JSON(http.StatusOK, resp)

	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "秒传成功",
	}
	c.JSON(http.StatusOK, resp)

}

// 返回文件下载地址
func DownloadUrlHandler(c *gin.Context) {
	// 获取文件hash
	filehash := c.Request.FormValue("filehash")

	//从tbl_file表中获取文件的信息
	row, err := dblayer.GetFmDb(filehash)
	if err != nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "文件不存在",
		}
		c.JSON(http.StatusOK, resp)

		return
	}

	// 判断文件存放是在本地还是在oss上

	// 本地

	// oss上
	url := store.DownloadUrl(row.FileHash, row.FileName.String)
	c.Data(http.StatusOK, "", []byte(url))

}

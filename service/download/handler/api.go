package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ggvylf/filestore/common"
	cfg "github.com/ggvylf/filestore/config"
	dbcli "github.com/ggvylf/filestore/service/dbproxy/client"
	store "github.com/ggvylf/filestore/store/minio"
)

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")

	// 从文件表查找记录
	dbResp, err := dbcli.GetFileMeta(filehash)
	if err != nil {
		log.Println(err)
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": common.StatusServerError,
				"msg":  "server error",
			})
		return
	}

	// 填充实例
	tblFile := dbcli.ToTableFile(dbResp.Data)

	// TODO: 判断文件存在OSS，还是Ceph，还是在本地
	// 本地的前缀是/tmp
	// 没用到ceph
	// oss是userfile/minio、
	if strings.HasPrefix(tblFile.FileAddr.String, cfg.TempLocalRootDir) ||
		strings.HasPrefix(tblFile.FileAddr.String, cfg.CephRootDir) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		// 调用/file/download接口
		tmpURL := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			c.Request.Host, filehash, username, token)
		c.Data(http.StatusOK, "application/octet-stream", []byte(tmpURL))
	} else if strings.HasPrefix(tblFile.FileAddr.String, cfg.OSSRootDir) {
		// oss下载url
		signedURL := store.DownloadUrl(tblFile.FileHash, tblFile.FileName.String)
		log.Println(tblFile.FileAddr.String)
		c.Data(http.StatusOK, "application/octet-stream", []byte(signedURL))
	}
}

// DownloadHandler : 文件下载接口
func DownloadHandler(c *gin.Context) {
	fsha1 := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")
	// TODO: 处理异常情况
	fResp, ferr := dbcli.GetFileMeta(fsha1)
	ufResp, uferr := dbcli.QueryUserFileMeta(username, fsha1)
	if ferr != nil || uferr != nil || !fResp.Suc || !ufResp.Suc {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": common.StatusServerError,
				"msg":  "server error",
			})
		return
	}
	uniqFile := dbcli.ToTableFile(fResp.Data)
	userFile := dbcli.ToTableUserFile(ufResp.Data)

	// 针对文件不同位置来处理

	if strings.HasPrefix(uniqFile.FileAddr.String, cfg.TempLocalRootDir) {
		// 本地文件， 直接下载
		c.FileAttachment(uniqFile.FileAddr.String, userFile.FileName)

		// ceph或者是其他存储
	} else if strings.HasPrefix(uniqFile.FileAddr.String, cfg.CephRootDir) {
		// 	// ceph中的文件，通过ceph api先下载
		// 	bucket := store.GetMC().GetCephBucket("userfile")
		// 	data, _ := bucket.Get(uniqFile.FileAddr.String)
		// 	//	c.Header("content-type", "application/octect-stream")
		// 	c.Header("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
		// 	c.Data(http.StatusOK, "application/octect-stream", data)
		// }

		//oss上，这里用的是minio
	} else if strings.HasPrefix(uniqFile.FileAddr.String, cfg.OSSBucket) {
		url := store.DownloadUrl(uniqFile.FileHash, userFile.FileName)
		c.Data(http.StatusOK, "", []byte(url))
	}
}

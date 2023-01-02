package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ggvylf/filestore/meta"
	"github.com/ggvylf/filestore/util"
)

// 上传文件并保存到本地
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		// 返回上传页面
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "InternalServerError")
			return
		}
		io.WriteString(w, string(data))

	} else if r.Method == "POST" {
		// 接收上传的内容并存储到本地

		//从form表单中获取文件
		file, head, err := r.FormFile("file")

		if err != nil {
			fmt.Printf("failed to get data,err=%v\n", err.Error())
			return
		}
		defer file.Close()

		// 初始化FileMeta
		fm := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UpoadAt:  time.Now().Format("2006-01-02 15:04:05"),
		}

		//新建一个本地文件的fd
		newfile, err := os.Create(fm.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err=%v\n", err.Error())
			return
		}
		defer newfile.Close()

		// 复制文件 同时可以获取文件大小
		fm.FileSize, err = io.Copy(newfile, file)
		if err != nil {
			fmt.Printf("Failed to write file,err=%v\n", err.Error())
			return
		}

		// 重置newfile的偏移量到文件头部
		newfile.Seek(0, 0)
		// 计算上传文件的sha1
		fm.FileSha1 = util.FileSha1(newfile)

		// append到元信息队列中
		meta.UploadFmList(fm)

		// 302重定向到上传成功页面
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传文件成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Success!")
}

// 返回元信息列表
func GetFmListHandler(w http.ResponseWriter, r *http.Request) {
	fmList := meta.GetFmList()
	data, err := json.Marshal(fmList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

// 返回指定sha1的fm对象
func GetFileMetaHander(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// filehash:=r.Form["filehash"][0]
	filehash := r.Form.Get("filehash")
	fm := meta.GetFm(filehash)
	data, err := json.Marshal(fm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

func DownFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	fm := meta.GetFm(filehash)

	//小文件可以 大文件性能不行 这两种方法是等价的
	//data, err := os.ReadFile(fm.Location)
	//data, err := ioutil.ReadFile(fm.FileName)

	// 先打开文件句柄再读取
	fd, err := os.Open(fm.FileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect")
	// 避免中文文件名乱码
	w.Header().Set("Content-Dispositon", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(fm.FileName)))

	w.Write(data)

}

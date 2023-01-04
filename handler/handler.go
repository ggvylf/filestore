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
		// meta.UploadFmList(fm)

		// 更新元数据到db
		_ = meta.UpdateFmDb(fm)

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
	// fm := meta.GetFm(filehash)

	// 从db中获取fm
	fm, err := meta.GetFmDb(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

// 文件下载
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

// 更新fm 元数据
func FmUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form.Get("filehash")
	opType := r.Form.Get("op")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	newFm := meta.GetFm(fileHash)
	newFm.FileName = newFileName
	meta.UploadFmList(newFm)

	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(newFm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// 删除fm和文件
func FmDeleteHander(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form.Get("filehash")

	fm := meta.GetFm(fileHash)

	// 这里注意要先删除文件 再删除元信息
	ok := meta.DeleteFile(fm.Location)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("file not exists"))
		return
	}

	meta.DeleteFm(fileHash)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete ok"))

}

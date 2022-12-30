package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

		//从from中获取文件
		file, head, err := r.FormFile("file")

		if err != nil {
			fmt.Printf("failed to get data,err=%v\n", err.Error())
			return
		}
		defer file.Close()

		//新建一个本地文件的fd
		newfile, err := os.Create("/tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("Failed to create file,err=%v\n", err.Error())
			return
		}
		defer newfile.Close()

		//复制文件
		_, err = io.Copy(newfile, file)
		if err != nil {
			fmt.Printf("Failed to write file,err=%v\n", err.Error())
			return
		}

		// 302重定向到上传成功页面
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传文件成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Success!")
}

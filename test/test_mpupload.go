package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	jsonit "github.com/json-iterator/go"
)

func multipartUpload(filename string, targetURL string, chunkSize int) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) //每次读取chunkSize大小的内容
	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}
		index++

		bufCopied := make([]byte, 5*1048576) // 5*1024*1024 5m
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			fmt.Printf("upload_size: %d\n", len(b))

			resp, err := http.Post(
				targetURL+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",
				bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}

			body, er := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v %+v\n", string(body), er)
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		//遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	for idx := 0; idx < index; idx++ {
		select {
		case res := <-ch:
			fmt.Println(res)
		}
	}

	return nil
}

func main() {
	username := "admin"
	token := "c439226a3f57e4f9ddc60a8e18d75d2963bbc5f5"
	filehash := "fe1d6ccb2544698b5c567411306e659de0fe922d" //sha1sum filename
	filesize := "148883574"
	filename := "/home/ggvylf/Downloads/go1.19.2.linux-amd64.tar.gz"

	// 1. 请求初始化分块上传接口
	resp, err := http.PostForm(
		"http://127.0.0.1:8888/file/mpupload/init",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {filesize},
		})

	if err != nil {
		fmt.Println("initmp resp err=", err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("initmp resp body read err=", err.Error())
		os.Exit(-1)
	}

	// 2. 得到uploadID以及服务端指定的分块大小chunkSize
	uploadID := jsonit.Get(body, "data").Get("UploadID").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	chunkCount := jsonit.Get(body, "data").Get("ChunkCount").ToInt()
	fmt.Printf("uploadid: %s  chunksize: %d chunkcount: %d\n", uploadID, chunkSize, chunkCount)

	// 3. 请求分块上传接口

	tURL := "http://127.0.0.1:8888/file/mpupload/uppart?" +
		"username=admin&token=" + token + "&uploadid=" + uploadID
	multipartUpload(filename, tURL, chunkSize)

	// 4. 请求分块完成接口
	resp, err = http.PostForm(
		"http://127.0.0.1:8888/file/mpupload/complete",
		url.Values{
			"username": {username},
			"token":    {token},
			"filehash": {filehash},
			"filesize": {filesize},
			"filename": {filename},
			"uploadid": {uploadID},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("complete result: %s\n", string(body))
}

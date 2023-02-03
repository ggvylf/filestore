package test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	username  = "admin"
	token     = "e922a114151039f67a9250bb9437772063dcae01"
	targetURL = "http://127.0.0.1:8888/file/fastupload"
	filehash  = "315d04cb9c699eb303fb8a39276d330737bbec70"
	filename  = "a.txt"
	filesize  = "29682"
)

func test_fastupload() {

	resp, err := http.PostForm(targetURL, url.Values{
		"username": {username},
		"token":    {token},
		"filehash": {filehash},
		"filename": {filename},
	})
	log.Printf("error: %+v\n", err)
	log.Printf("resp: %+v\n", resp)
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		log.Printf("parseBodyErr: %+v\n", err)
		if err == nil {
			log.Printf("parseBody: %+v\n", string(body))
		}
	}
}

func main() {
	test_fastupload()
}

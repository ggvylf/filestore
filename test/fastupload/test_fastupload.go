package main

import (
	"io"
	"log"
	"net/http"
	"net/url"

	config "github.com/ggvylf/filestore/test/config"
)

func test_fastupload() {

	resp, err := http.PostForm(config.FasttargetURL, url.Values{
		"username": {config.Username},
		"token":    {config.Token},
		"filehash": {config.Fastfilehash},
		"filename": {config.Fastfilename},
	})
	log.Printf("error: %+v\n", err)
	log.Printf("resp: %+v\n", resp)
	if resp != nil {
		body, err := io.ReadAll(resp.Body)

		log.Printf("parseBodyErr: %+v\n", err)
		if err == nil {
			log.Printf("parseBody: %+v\n", string(body))
		}
	}
}

func main() {
	test_fastupload()
}

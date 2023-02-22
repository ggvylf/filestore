package main

import (
	"github.com/ggvylf/filestore/service/apigw/route"
)

func main() {

	r := route.Router()
	r.Run(":8888")
}

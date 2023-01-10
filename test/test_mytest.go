package main

import (
	"fmt"
	"os"
)

func main() {

	dir := "/tmp/a/"
	fm := "a"

	fd, _ := os.OpenFile(dir+fm, os.O_CREATE|os.O_APPEND, 0644)
	defer fd.Close()

	files, _ := os.ReadDir(dir)
	for _, f := range files {
		data, err := os.ReadFile(dir + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		fd.Write(data)
	}
}

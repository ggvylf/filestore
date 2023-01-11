package main

func main() {

	// dir := "/tmp/a/"
	// fm := "a"
	// // os.OpenFile默认权限是os.O_RDONLY，需要添加os.O_WRONLY或者是os.O_RDWR
	// fd, _ := os.OpenFile(dir+fm, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	// defer fd.Close()

	// files, _ := os.ReadDir(dir)
	// for _, f := range files {
	// 	// Name()只返回文件名，不包含路径
	// 	if f.Name() == fm {
	// 		break
	// 	}
	// 	data, err := os.ReadFile(dir + f.Name())
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	_, err = fd.Write(data)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	// filename := "/tmp/a/a"
	// f := strings.Split(filename, "/")
	// fmt.Println(f[len(f)-1])

}

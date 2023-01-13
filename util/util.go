package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

// 文件名排序
// https://blog.azhangbaobao.cn/2020/07/09/go%E8%AF%AD%E8%A8%80%E5%AE%9E%E7%8E%B0%E5%A4%A7%E6%96%87%E4%BB%B6%E5%88%86%E5%89%B2%E4%B8%8E%E5%90%88%E5%B9%B6.html
// https://gosamples.dev/sort-string-numbers/

func FileSortForNumber(data []string) ([]string, error) {
	var lastErr error
	sort.Slice(data, func(i, j int) bool {
		a, err := strconv.ParseInt(filepath.Base(data[i]), 10, 64)
		if err != nil {
			lastErr = err
			return false
		}
		b, err := strconv.ParseInt(filepath.Base(data[j]), 10, 64)
		if err != nil {
			lastErr = err
			return false
		}

		// 定义排序规则 这里是从小到大
		return a < b
	})
	return data, lastErr
}

// 带前缀的文件名排序

func SplitName(file string) (string, error) {

	// file /tmp/admin1739c1065654ebf7/aa29.txt
	// base="aa29.txt"
	// ext="txt"
	// name="aa29"

	base := filepath.Base(file)
	ext := filepath.Ext(file)
	name := base[:len(base)-len(ext)]

	// 计算文件名中数字的个数
	i := len(name) - 1
	for ; i >= 0; i-- {
		if '0' > name[i] || name[i] > '9' {
			break
		}
	}
	i++

	// string numeric suffix to uint64 bytes
	// empty string is zero, so integers are plus one

	// 把文件名转换成uint类型
	b64 := make([]byte, 64/8)
	s64 := name[i:]

	if len(s64) > 0 {
		u64, err := strconv.ParseUint(s64, 10, 64)
		if err == nil {
			binary.BigEndian.PutUint64(b64, u64+1)
		}
	}

	// prefix + numeric-suffix + ext
	return name[:i] + string(b64) + ext, nil
}

func FileSortForStringWithNum(data []string) ([]string, error) {
	var lastErr error

	sort.Slice(data, func(i, j int) bool {

		a, err := SplitName(data[i])
		if err != nil {
			lastErr = err
			return false
		}
		b, err := SplitName(data[j])
		if err != nil {
			lastErr = err
			return false
		}

		return a < b

	})

	return data, lastErr
}

package helpers

import (
	"crypto/md5"
	"fmt"
)

// HashMd5ByString md5 哈希函数, 输出小写的十六进制
func HashMd5ByString(str string) string {
	data := []byte(str)
	hashStr := md5.Sum(data)
	md5str := fmt.Sprintf("%x", hashStr)
	return md5str
}

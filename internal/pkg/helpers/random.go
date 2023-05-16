package helpers

import (
	"math/rand"
	"sync"
	"time"
)

const (
	charset       = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
	numberCharset = "0123456789"
)

var seededRand *rand.Rand = rand.New(&safeRand{src: rand.NewSource(time.Now().UnixNano()).(rand.Source64)})

// RandomStringWithCharset 从 charset 中生成指定 length 的字符串
func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	lc := len(charset)
	for i := range b {
		b[i] = charset[seededRand.Intn(lc)]
	}
	return string(b)
}

// RandomString 生成 length 长度的字符串
//
//	生成种子使用 const charset
func RandomString(length int) string {
	return RandomStringWithCharset(length, charset)
}

// RandomNumberString 生成 length 长度的数字字符串
//
//	生成种子使用 const numberCharset
func RandomNumberString(length int) string {
	return RandomStringWithCharset(length, numberCharset)
}

// 并发安全的随机数种子
type safeRand struct {
	lk  sync.Mutex
	src rand.Source64
}

func (r *safeRand) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *safeRand) Uint64() (n uint64) {
	r.lk.Lock()
	n = r.src.Uint64()
	r.lk.Unlock()
	return
}

func (r *safeRand) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

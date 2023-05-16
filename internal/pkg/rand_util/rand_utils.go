package rand_util

import (
	"math/rand"
	"sync"
	"time"
)

var once sync.Once

func init() {
	once.Do(func() {
		rand.Seed(time.Now().Unix())
	})
}

func Get(min, max int) int {
	return rand.Intn(max-min) + min
}

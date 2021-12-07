package utils

import (
	"math/rand"
	"time"
)

// RandomNumber 指定位数随机数
func RandomNumber(lo, hi int) int {
	rand.Seed(time.Now().UnixNano())
	return lo + rand.Intn(hi-lo)
}

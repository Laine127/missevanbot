package utils

import (
	"math/rand"
	"time"
)

// RandomNumber generate a random number between lo and hi.
func RandomNumber(lo, hi int) int {
	rand.Seed(time.Now().UnixNano())
	return lo + rand.Intn(hi-lo)
}

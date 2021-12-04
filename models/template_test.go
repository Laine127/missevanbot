package models

import (
	"log"
	"testing"
	"time"
)

func TestLuckString(t *testing.T) {
	for i := 0; i < 100; i++ {
		log.Println(LuckString())
		time.Sleep(time.Millisecond)
	}
}
